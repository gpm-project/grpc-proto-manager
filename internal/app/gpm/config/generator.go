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

package config

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

// GeneratorType defining the enum with proto generators.
type GeneratorType int

const (
	// DockerCmd proto generator.
	DockerCmd GeneratorType = iota
	// DockerizedCmd to use the embeeded proto generator.
	DockerizedCmd
)

// GeneratorTypeToString map associating type an string representation.
var GeneratorTypeToString = map[GeneratorType]string{
	DockerCmd:     "docker",
	DockerizedCmd: "dockerized",
}

// GeneratorTypeToEnum map associating string representation with enum type.
var GeneratorTypeToEnum = map[string]GeneratorType{
	"docker":     DockerCmd,
	"dockerized": DockerizedCmd,
}

type GeneratorConfig struct {
	// GeneratorName with the name of the provider implementing the operations of proto code generation.
	GeneratorName string
	// DockerCmdImage with the image to be used to generate the protos when using an auxiliary docker command.
	DockerCmdImage string
}

// IsValid checks if the configuration options are valid.
func (gc *GeneratorConfig) IsValid() error {
	if gc.GeneratorName != GeneratorTypeToString[DockerCmd] && gc.GeneratorName != GeneratorTypeToString[DockerizedCmd] {
		return fmt.Errorf("generatorName cannot be empty, valid values: docker, dockerized")
	}
	return nil
}

// Print the configuration using the application logger.
func (gc *GeneratorConfig) Print() {
	// Use logger to print the configuration
	log.Info().Str("generator", gc.GeneratorName).Msg("Providers")
}
