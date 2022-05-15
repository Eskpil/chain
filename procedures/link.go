package procedures

import "chain/compilers"

type LinkProcedure struct {
	Files  []string
	Into   string
	Linker compilers.Compiler
}

func (p LinkProcedure) RunProcedure() error {
	return nil
}
