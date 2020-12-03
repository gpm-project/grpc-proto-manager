package main

import (
	"os"

	"github.com/gpm-project/grpc-proto-manager/cmd/gpm/commands"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Version of the command
var Version string

// Commit from which the command was built
var Commit string

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	commands.Execute(Version, Commit)
}
