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

package protos

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
)

// DockerizedCmdProvider is a proto generator based on issuing commands inside a docker container.
// This command expects to be executed inside the namely/protoc-all docker image.
type DockerizedCmdProvider struct {
	Common
}

// NewDockerizedCmdGenerator assumes the proto tools are embeeded and locally available.
func NewDockerizedCmdGenerator() (Generator, error) {
	log.Debug().Msg("Using DockerizedCmd proto generator")
	return &DockerizedCmdProvider{Common: Common{}}, nil
}

// Generate a set of proto stubs in a given language.
func (dcp *DockerizedCmdProvider) Generate(rootPath string, targetName string, generatedPath string, language string) error {
	// Based on the documentation available at: https://github.com/namely/docker-protoc
	log.Debug().Str("rootPath", rootPath).Str("targetName", targetName).Str("generatedPath", generatedPath).Str("language", language).Msg("generating protos")

	cmdArgs := []string{
		"-l", language, // Target language
		"-d", targetName, // Directory to take protos from
		"-i", ".", // Include local path
		"-o", "generated", // Path where the resulting code is stored.
	}

	// Extra options from the available arguments
	if language == "go" {
		cmdArgs = append(cmdArgs, "--with-gateway") //Generate grpc-gateway files (experimental)
	}

	if language == "go" || language == "gogo" || language == "cpp" || language == "java" || language == "python" {
		cmdArgs = append(cmdArgs, "--with-validator") // Generate validations for (go gogo cpp java python)
	}

	cmd := exec.Command("entrypoint.sh", cmdArgs...)
	log.Debug().Interface("cmd", cmd.Args).Msg("dockerized generation cmd")
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
