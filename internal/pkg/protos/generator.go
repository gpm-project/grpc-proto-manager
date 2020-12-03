package protos

import "fmt"

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

// Generator interface for all implementations.
type Generator interface {
	// Generate a set of proto stubs in a given language.
	Generate(rootPath string, targetName string, generatedPath string, language string) error
}

// NewGenerator builds a new generator.
func NewGenerator(generatorName string) (Generator, error) {
	gen, exists := GeneratorTypeToEnum[generatorName]
	if !exists {
		return nil, fmt.Errorf("generator %s not found", generatorName)
	}
	switch gen {
	case DockerCmd:
		return NewDockerCmdGenerator()
	case DockerizedCmd:
		return NewDockerizedCmdGenerator()
	}
	return nil, fmt.Errorf("no implementation found for %s generator", generatorName)
}
