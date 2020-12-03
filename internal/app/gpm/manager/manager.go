package manager

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/gpm-project/grpc-proto-manager/internal/app/gpm/config"
	"github.com/gpm-project/grpc-proto-manager/internal/pkg/files"
	"github.com/gpm-project/grpc-proto-manager/internal/pkg/protos"
	"github.com/gpm-project/grpc-proto-manager/internal/pkg/repo"
	"github.com/rs/zerolog/log"
)

// ProtoLangFileName defines the name of the file that specifies the target languages.
const ProtoLangFileName = ".protolangs"

// ExcludedDirs with the list of directories that will be excluded by default.
var ExcludedDirs = []string{".git", ".github"}

// GPM structure with the manager main loop.
type GPM struct {
	cfg                config.ServiceConfig
	repositoryProvider repo.Provider
	protoGenerator     protos.Generator
}

// NewManager creates a new GPM entity.
func NewManager(cfg config.ServiceConfig) *GPM {
	return &GPM{cfg: cfg}
}

// SetupGeneratorConfig determines if the required parameters are present depending on the selected environment.
func (gpm *GPM) SetupGeneratorConfig() error {
	generator, exists := protos.GeneratorTypeToEnum[gpm.cfg.GeneratorName]
	if !exists {
		return fmt.Errorf("unsupported generator %s", gpm.cfg.GeneratorName)
	}

	repoProvider, exists := repo.RepositoryTypeToEnum[gpm.cfg.RepositoryProvider]
	if !exists {
		return fmt.Errorf("unsupported repository provider %s", gpm.cfg.RepositoryProvider)
	}

	switch generator {
	case protos.DockerizedCmd:
		return gpm.SetupDockerizedGeneration(repoProvider)
	}
	return nil
}

// SetupDockerizedGeneration configures the different provider attending to the execution environment. In this
// case, as we are being executed from within a docker container, git information needs to be set.
func (gpm *GPM) SetupDockerizedGeneration(repoProvider repo.RepositoryType) error {
	if gpm.cfg.RepositoryPusherUsername == "" {
		return fmt.Errorf("--repositoryPusherUsername is required when running in a containerized environment")
	}
	if gpm.cfg.RepositoryPusherEmail == "" {
		return fmt.Errorf("--repositoryPusherEmail is required when running in a containerized environment")
	}
	return gpm.repositoryProvider.ConfigurePusher(gpm.cfg.RepositoryPusherUsername, gpm.cfg.RepositoryPusherEmail, gpm.cfg.RepositoryAccessToken)
}

// Run triggers the execution of the command.
func (gpm *GPM) Run(basePath string) error {
	log.Debug().Msg("Launching GPM")
	if err := gpm.cfg.IsValid(); err != nil {
		log.Fatal().Err(err).Msg("invalid configuration options")
	}
	gpm.cfg.Print()
	defer gpm.cleanup(basePath, gpm.cfg.TempPath)

	repoProvider, err := repo.NewRepoProvider(gpm.cfg.RepositoryProvider)
	if err != nil {
		return err
	}
	gpm.repositoryProvider = repoProvider

	protoGenerator, err := protos.NewGenerator(gpm.cfg.GeneratorName)
	if err != nil {
		return err
	}
	gpm.protoGenerator = protoGenerator

	err = gpm.SetupGeneratorConfig()
	if err != nil {
		return err
	}

	// Iterate over the project directories

	fileInfo, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, info := range fileInfo {
		if info.IsDir() && !gpm.isExcluded(info.Name()) {
			targetPath := path.Join(basePath, info.Name())
			err := gpm.ProcessProtoDirectory(targetPath, info.Name())
			if err != nil {
				return err
			}
		}
	}

	return err
}

// cleanup function to delete temporal directories.
func (gpm *GPM) cleanup(basePath string, tempPath string) {
	log.Debug().Str("basePath", basePath).Str("tempPath", tempPath).Msg("cleaning temporal directories")
	generatedPath := path.Join(basePath, "generated")
	_, err := os.Stat(generatedPath)
	if err == nil {
		// Otherwise, assume directory has not being generated.
		if rerr := os.RemoveAll(generatedPath); rerr != nil {
			log.Warn().Str("generatedPath", generatedPath).Err(rerr).Msg("unable to deleted directory with temporal generated code")
		}
	}
	_, err = os.Stat(tempPath)
	if err == nil {
		// Otherwise, assume directory has not being generated.
		if rerr := os.RemoveAll(tempPath); rerr != nil {
			log.Warn().Str("tempPath", tempPath).Err(rerr).Msg("unable to deleted directory with temporal downloaded code")
		}
	}
}

// isExcluded method to check if the directory should be excluded from the generation process.
func (gpm *GPM) isExcluded(dirName string) bool {
	for _, excludedDir := range ExcludedDirs {
		if excludedDir == dirName {
			return true
		}
	}
	return false
}

// LoadProtoLangs loads the file in each directory that defines the target languages. If none is found, the default one for
// the project will be returned.
func (gpm *GPM) LoadProtoLangs(targetPath string) ([]string, error) {
	protolangsFile := path.Join(targetPath, ProtoLangFileName)
	if _, err := os.Stat(protolangsFile); os.IsNotExist(err) {
		return []string{gpm.cfg.DefaultLanguage}, nil
	}
	// read the file, no long lines are expected
	readFile, err := os.Open(targetPath)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)
	langs := make([]string, 0)
	for scanner.Scan() {
		langs = append(langs, scanner.Text())
	}
	readFile.Close()
	return langs, nil
}

// getRepoName obtains the name of the target repository associated with a given directory and language
func (gpm *GPM) getRepoName(directoryName string, language string) string {
	return fmt.Sprintf("grpc-%s-%s", directoryName, language)
}

// ProcessProtoDirectory is the main function to compile, calculate the difference in code with the previous version, and commit the changes.
func (gpm *GPM) ProcessProtoDirectory(targetPath string, name string) error {
	log.Info().Str("path", targetPath).Msg("processing proto directory")

	targetLanguages, err := gpm.LoadProtoLangs(targetPath)
	if err != nil {
		return err
	}
	log.Debug().Interface("languages", targetLanguages).Msg("target")
	for _, language := range targetLanguages {
		repoName := gpm.getRepoName(name, language)
		repoURL, err := gpm.repositoryProvider.GetRepoURL(gpm.cfg.RepositoryOrganization, repoName)
		if err != nil {
			return fmt.Errorf("cannot determine repository URL: %w", err)
		}
		// First step is to clone the generated proto repo to compare the files. Notice that generated files have timestamped data,
		// and diff is not recommended on that data.
		tmpRepoDir := path.Join(gpm.cfg.TempPath, repoName)
		err = gpm.repositoryProvider.Clone(repoURL, tmpRepoDir)
		if err != nil {
			return fmt.Errorf("cannot clone target repository %s to calculate diff: %w", repoURL, err)
		}

		// Now compare the content.
		equal, err := files.CompareDirectoriesAreEqual(".proto", targetPath, tmpRepoDir)
		if err != nil {
			return fmt.Errorf("cannot compare files: %w", err)
		}
		if !equal {
			// If there is a change, generate the proto stubs on the given languages
			err := gpm.OrchestrateGeneration(name, tmpRepoDir, language)
			if err != nil {
				return fmt.Errorf("cannot generate proto code: %w", err)
			}
		} else {
			log.Info().Str("repo", repoName).Msg("no changes detected, skipping generation")
		}
		// Remove the temporal directory
		_ = os.RemoveAll(tmpRepoDir)
	}

	return nil
}

// OrchestrateGeneration orchestrates the generation of the protos.
func (gpm *GPM) OrchestrateGeneration(name string, tmpRepoDir string, language string) error {
	// Generate the code
	err := gpm.protoGenerator.Generate(gpm.cfg.ProjectPath, name, tmpRepoDir, language)
	if err != nil {
		return err
	}
	// Calculate version
	version, err := gpm.repositoryProvider.GetLastVersion(tmpRepoDir)
	if err != nil {
		return err
	}
	log.Debug().Str("previous", version.String()).Msg("version")
	// Publish if required
	if gpm.cfg.SkipPublish {
		log.Warn().Str("repo", tmpRepoDir).Msg("changes will not be published")
		return nil
	}
	version.IncrementMinor()
	log.Info().Str("newVersion", version.String()).Str("repo", gpm.getRepoName(name, language)).Msg("publishing new version")
	return gpm.repositoryProvider.Publish(tmpRepoDir, version)
}
