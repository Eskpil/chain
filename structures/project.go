package structures

import (
	"chain/logger"
	"os"
)

type ProjectStructure struct {
	Project struct {
		Name        *string
		Procedures  *[]string
		Compilers   *string
		SubProjects []string
	}
}

func (project *ProjectStructure) Validate(path string) {
	logger.Info.Printf("Validating project: %s\n", path)

	if project.Project.Name == nil {
		logger.Error.Printf("%s: Project names are mandatory\n", path)
		os.Exit(1)
	}

	if project.Project.Procedures == nil {
		logger.Warn.Printf("%s: Project has no procedures\n", path)
	}
}
