package structures

type ProcedureStructure struct {
	Procedure struct {
		Name  string
		Build *struct {
			Files []string
		}
		Link *struct {
			Files  []string
			Target string
			Into   string
			With   []string
		}
		Library *struct {
			Name string
			From string
		}
		Export []string
	}
}
