package commands

import (
	"github.com/gpm-project/grpc-proto-manager/internal/app/gpm/manager"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generateCmdLongHelp = `
This command triggers the generation of the proto stubs for a collection of protos.
`

var generateCmdExamples = `
# Generate all the protos from the current directory.
$ gpm generate .
`

var generateCmd = &cobra.Command{
	Use:     "generate <base_path>",
	Short:   "Generate the resulting stubs for a collection of proto specs",
	Long:    generateCmdLongHelp,
	Example: generateCmdExamples,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		readConfig(args[0])
		gpm := manager.NewManager(appConfig)
		err := gpm.Run(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("generation failed")
		}
	},
}

func init() {
	generateCmd.Flags().String("tempPath", "/tmp/gpm",
		"Temporal file for the generation of intermediate data")
	generateCmd.Flags().StringVar(&appConfig.GeneratorName, "protoGenerator", "docker", "Implementation used to generate the proto code.")
	generateCmd.Flags().StringVar(&appConfig.RepositoryAccessToken, "repositoryAccessToken", "", "An access token for the authentication of the repository provider. Use this for GitHub actions.")
	generateCmd.Flags().BoolVar(&appConfig.SkipPublish, "skipPublish", false, "Flag to skip publishing the generated protos")
	err := viper.BindPFlag("tempPath", generateCmd.Flags().Lookup("tempPath"))
	if err != nil {
		log.Error().Err(err).Msg("unable to bind viper key")
	}

	rootCmd.AddCommand(generateCmd)
}

// readConfig gets the project configuration and applies it.
func readConfig(fromPath string) {
	viper.SetEnvPrefix("GPM")
	viper.AutomaticEnv()
	viper.AddConfigPath(fromPath)
	viper.SetConfigName(".gpm") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn().Msg("No config file found on given path, create a .gpm.yaml file for consistent results.")
		} else {
			log.Fatal().Err(err).Msg("unable to read configuration file")
		}
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatal().Err(err).Msg("unable to unmarshal resolved configuration into config structure. Check structure/file structure for a mismatch")
	}
	appConfig.ProjectPath = fromPath
	log.Info().Str("path", viper.ConfigFileUsed()).Msg("configuration loaded")
}
