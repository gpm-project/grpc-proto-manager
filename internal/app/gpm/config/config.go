package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ServiceConfig structure with all the options required by the service and service components.
type ServiceConfig struct {
	// Debug level activated.
	Debug bool
	// Version of the grpc-proto-manager tool.
	Version string
	// Commit with the built of the grpc-proto-manager tool.
	Commit string
	// RepositoryProvider with the target repository provider (e.g., GitHub).
	RepositoryProvider string
	// RepositoryOrganization with the organization that contains the generated code.
	RepositoryOrganization string
	// RepositoryUsername with the name of the actor pushing the changes. This value is required if GPM is executed from within a container.
	RepositoryPusherUsername string
	// RepositoryEmail with the email pushing the changes. This value is required if GPM is executed from within a container.
	RepositoryPusherEmail string
	// RepositoryAccessToken with a token required to access the repository. This value is required for the github action provider.
	RepositoryAccessToken string
	// DefaultLanguage to generate the protos if not .protolangs file is found.
	DefaultLanguage string
	// ProjectPath with the path of the gRPC proto repo being analyzed.
	ProjectPath string
	// TempPath with the path used to generated temporal data.
	TempPath string
	// SkipPublish determines if the generated protos are to be published.
	SkipPublish bool
	// GeneratorName with the name of the provider implementing the operations of proto code generation.
	GeneratorName string
}

// resolvePath processes the path given as input and translates it based on relative abstractions.
func (sc *ServiceConfig) resolvePath(targetPath string) string {
	// TODO Process ., .., ../., ~, etc.
	return targetPath
}

// createDirectoryIfNotExists checks if a directory exists and creates it otherwise.
func (sc *ServiceConfig) createDirectoryIfNotExists(targetPath string) error {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err = os.MkdirAll(targetPath, 0755)
		if err != nil {
			return fmt.Errorf("Unable to create directory %w", err)
		}
	}
	return nil
}

// logSetup sets the desired log level.
func (sc *ServiceConfig) logSetup() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if sc.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// IsValid checks if the configuration options are valid.
func (sc *ServiceConfig) IsValid() error {
	sc.logSetup()
	if sc.ProjectPath == "" {
		return fmt.Errorf("projectPath cannot be empty")
	}
	if sc.RepositoryProvider == "" {
		return fmt.Errorf("repositoryProvider cannot be empty")
	}
	if sc.RepositoryOrganization == "" {
		return fmt.Errorf("repositoryOrganization cannot be empty")
	}
	if sc.DefaultLanguage == "" {
		return fmt.Errorf("defaultLanguage cannot be empty")
	}
	if err := sc.createDirectoryIfNotExists(sc.TempPath); err != nil {
		return err
	}
	return nil
}

// Print the configuration using the application logger.
func (sc *ServiceConfig) Print() {
	// Use logger to print the configuration
	log.Info().Str("version", sc.Version).Str("commit", sc.Commit).Msg("app config")
	log.Info().Str("Project", sc.ProjectPath).Str("Temp", sc.TempPath).Msg("Paths")
	log.Info().Str("Repository", sc.RepositoryProvider).Str("generator", sc.GeneratorName).Msg("Providers")
	log.Info().Str("Language", sc.DefaultLanguage).Msg("Defaults")
	if sc.SkipPublish {
		log.Warn().Msg("Proto publication is disabled")
	}
	log.Info().Str("URL", sc.RepositoryOrganization).Msg("generated code repository")
	// Pusher related information.
	pusherInfo := log.Info()
	if sc.RepositoryPusherUsername == "" {
		pusherInfo = pusherInfo.Str("username", "<system default>")
	} else {
		pusherInfo = pusherInfo.Str("username", sc.RepositoryPusherUsername)
	}
	if sc.RepositoryPusherEmail == "" {
		pusherInfo = pusherInfo.Str("email", "<system default>")
	} else {
		pusherInfo = pusherInfo.Str("email", sc.RepositoryPusherEmail)
	}
	if sc.RepositoryAccessToken != "" {
		pusherInfo = pusherInfo.Str("accessToken", strings.Repeat("*", len(sc.RepositoryAccessToken)))
	}

	pusherInfo.Msg("pusher information")
}
