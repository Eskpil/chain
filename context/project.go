package context

import (
	"chain/logger"
	"chain/structures"
	"chain/util"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func (s *Scope) RunProject(project *structures.ProjectStructure) {
	compilers := make(map[string]structures.Compiler)

	for _, c := range util.LoadDefaultCompilers().Compilers {
		compilers[c.Name] = c
	}

	if project.Project.Compilers != nil {
		structure := structures.LoadCompilersFrom(*project.Project.Compilers)

		for _, compiler := range structure.Compilers {
			compilers[compiler.Name] = compiler
		}
	}

	if s.Root {
		s.LoadHooks()
	}

	s.Compilers = compilers

	logger.Info.Printf("Running all procedures for project: %s\n", *project.Project.Name)

	if len(project.Project.Procedures) > 0 {
		for _, proc := range project.Project.Procedures {
			filepath := fmt.Sprintf("%s/%s/procedure.yml", s.Prefix, proc)

			data, err := os.ReadFile(filepath)

			if err != nil {
				logger.Error.Printf("Unexpected error when reading data\n")
				logger.PrintError(fmt.Sprintf("%s\n", err))
				os.Exit(1)
			}

			procedure := structures.ProcedureStructure{}

			err = yaml.Unmarshal(data, &procedure)

			if err != nil {
				logger.Error.Printf("Unexpected error when unmarshaling data\n")
				logger.PrintError(fmt.Sprintf("%s\n", err))
				os.Exit(1)
			}

			procedure.Validate(filepath)

			childScope := Scope{}

			childScope.InheritFrom(s, proc)

			childScope.RunProcedure(procedure)

			childScope.ExportUpwards()
		}
	}

	if len(project.Project.SubProjects) > 0 {
		for _, path := range project.Project.SubProjects {
			filepath := fmt.Sprintf("%s/%s/project.yml", s.Prefix, path)

			data, err := os.ReadFile(filepath)

			subproject := structures.ProjectStructure{}

			err = yaml.Unmarshal(data, &subproject)

			if err != nil {
				logger.Error.Printf("Unexpected error when unmarshaling data\n")
				logger.PrintError(fmt.Sprintf("%s\n", err))
				os.Exit(1)
			}

			project.Validate(*subproject.Project.Name)

			childScope := Scope{}

			childScope.InheritFrom(s, path)

			childScope.RunProject(&subproject)

			childScope.ExportUpwards()
		}
	}
}
