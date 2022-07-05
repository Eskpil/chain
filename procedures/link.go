package procedures

import (
	"chain/compilers"
)

type Target int64

const (
	Binary  Target = iota
	Library Target = iota
)

type LinkProcedure struct {
	Files  []string
	Into   string
	With   []compilers.Library
	Target Target
	Linker compilers.Compiler
}

func (p LinkProcedure) RunProcedure() error {
	var err error
	if p.Target == Library {
		err = p.Linker.LinkLibrary(p.Files, p.Into, p.With)
	} else {
		err = p.Linker.LinkBinary(p.Files, p.Into, p.With)
	}

	if err != nil {
		return err
	}

	return nil
}
