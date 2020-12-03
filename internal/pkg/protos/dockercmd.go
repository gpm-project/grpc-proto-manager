package protos

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
)

// DockerCmdProvider is a proto generator based on issuing docker commands. Future
// implementations will rely on the docker library.
type DockerCmdProvider struct {
	Common
}

// NewDockerCmdGenerator uses an external command to launch the docker container with the proto tools.
func NewDockerCmdGenerator() (Generator, error) {
	log.Debug().Msg("Using DockerCmd proto generator")
	return &DockerCmdProvider{Common: Common{}}, nil
}

// Generate a set of proto stubs in a given language.
func (dcp *DockerCmdProvider) Generate(rootPath string, targetName string, generatedPath string, language string) error {
	// Based on the documentation available at: https://github.com/namely/docker-protoc

	log.Debug().Str("rootPath", rootPath).Str("targetName", targetName).Str("generatedPath", generatedPath).Str("language", language).Msg("generating protos")

	cmdArgs := []string{
		"run",
		"-v", fmt.Sprintf("%s:/defs", rootPath), // source proto definition. This should be the root so imports work :)
		"namely/protoc-all:1.32_4", // Image, maybe move this as a constant or config value.
		"-l", language,             // Target language
		"-d", targetName, // Directory to take protos from
		"-i", ".", // Include local path
		"-o", "generated", // Path where the resulting code is stored.
		// Extra options from the available arguments
		"--with-gateway",   //Generate grpc-gateway files (experimental)
		"--with-validator", // Generate validations for (go gogo cpp java python)
	}

	cmd := exec.Command("docker", cmdArgs...)
	log.Debug().Interface("cmd", cmd).Msg("docker generation cmd")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to generate protos for %s due to %w: %s", targetName, err, string(stdoutStderr))
	}
	log.Debug().Msg("execution finished")
	log.Debug().Str("output", string(stdoutStderr)).Msg("protos successfully generated")
	err = dcp.copyAllSourceFiles(path.Join(rootPath, targetName), generatedPath)
	if err != nil {
		return fmt.Errorf("unable to copy source files: %w", err)
	}
	return dcp.moveGeneratedFiles(path.Join(rootPath, "generated"), generatedPath)
}
