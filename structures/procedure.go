package structures

type ProcedureStructure struct {
	Procedure struct {
		Name  string
		Build struct {
			Files []string
		}
		Link struct {
			Files  []string
			Target string
			Into   string
		}
	}
}
