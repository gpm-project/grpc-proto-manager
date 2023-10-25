/**
 * Copyright 2023 GPM Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
