package structures

import (
	"chain/logger"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Compiler struct {
	Name     string
	Path     string
	Language string
	Flags    []string
}

type CompilersStructure struct {
	Compilers []Compiler
}

func LoadCompilersFrom(path string) CompilersStructure {
	structure := CompilersStructure{}

	data, err := os.ReadFile(path)

	if err != nil {
		logger.Error.Printf("Unexpected error when reading data\n")
		logger.PrintError(fmt.Sprintf("%s\n", err))
		os.Exit(1)
	}

	err = yaml.Unmarshal(data, &structure)

	structure.Validate()

	return structure
}

func (structure *CompilersStructure) Validate() {
	for _, compiler := range structure.Compilers {
		if 0 >= len(compiler.Name) {
			logger.Error.Printf("Compiler.Name must not be empty.\n")
			os.Exit(1)
		}

		if 0 >= len(compiler.Language) {
			logger.Error.Printf("Compiler.Language must not be empty.\n")
			os.Exit(1)
		}

		if 0 >= len(compiler.Path) {
			logger.Error.Printf("Compiler.Path must not be emtpy.\n")
			os.Exit(1)
		}
	}
}
