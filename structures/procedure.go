package structures

type ProcedureStructure struct {
	Procedure struct {
		Name  string
		Build *struct {
			Compiler string
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
