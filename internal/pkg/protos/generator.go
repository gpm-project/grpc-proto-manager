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
	"github.com/gpm-project/grpc-proto-manager/internal/app/gpm/config"
)

// Generator interface for all implementations.
type Generator interface {
	// Generate a set of proto stubs in a given language.
	Generate(rootPath string, targetName string, generatedPath string, language string) error
}

// NewGenerator builds a new generator.
func NewGenerator(generatorConfig *config.GeneratorConfig) (Generator, error) {
	gen, exists := config.GeneratorTypeToEnum[generatorConfig.GeneratorName]
	if !exists {
		return nil, fmt.Errorf("generator %s not found", generatorConfig.GeneratorName)
	}
	switch gen {
	case config.DockerCmd:
		return NewDockerCmdGenerator(generatorConfig)
	case config.DockerizedCmd:
		return NewDockerizedCmdGenerator()
	}
	return nil, fmt.Errorf("no implementation found for %s generator", generatorConfig.GeneratorName)
}
