package procedures

import "chain/compilers"

type BuildProcedure struct {
	Files    []string
	Compiler compilers.Compiler
}

func (p BuildProcedure) RunProcedure() error {
	return nil
}
