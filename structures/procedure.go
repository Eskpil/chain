package structures

import (
	"chain/logger"
	"os"
)

type ProcedureStructure struct {
	Procedure *struct {
		Name  *string
		Build *struct {
			Hook     []string
			Compiler string
			Headers  *string
			Files    []string
		}
		Link *struct {
			Files  []string
			Target string
			Into   string
			Linker string
			With   []struct {
				Name string
				Kind string
			}
		}
		Library *struct {
			Name string
			From string
		}
		Export []string
	}
}

func (proc *ProcedureStructure) Validate(path string) {
	logger.Info.Printf("Validating procedure: %s\n", path)
	if proc.Procedure == nil {
		logger.Error.Printf("%s: Expected procedure information found nothing\n", path)
		os.Exit(1)
	}

	if proc.Procedure.Name == nil {
		logger.Error.Printf("%s: Expected procedure name\n", path)
		os.Exit(1)
	}

	if proc.Procedure.Build == nil {
		logger.Error.Printf("%s: Subprocedure Build is mandatory\n", path)
		os.Exit(1)
	}
}
