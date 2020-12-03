package commands

import (
	"fmt"
	"os"

	"github.com/gpm-project/grpc-proto-manager/internal/app/gpm/config"
	"github.com/spf13/cobra"
)

var appConfig config.ServiceConfig

var rootCmd = &cobra.Command{
	Use:     "gpm",
	Short:   "gRPC proto manager",
	Long:    `A simple manager to orchestrate the generation of gRPC protos`,
	Version: "TBD",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute the user command
func Execute(version string, commit string) {
	versionTemplate := fmt.Sprintf("version: %s commit: %s", version, commit)
	rootCmd.SetVersionTemplate(versionTemplate)
	appConfig.Version = version
	appConfig.Commit = commit
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&appConfig.Debug, "debug", false, "Enable debug log")
}
