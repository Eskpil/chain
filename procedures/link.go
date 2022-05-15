package procedures

import (
	"chain/compilers"
	"os"
)

type Target int64

const (
	Binary  Target = iota
	Library Target = iota
)

type LinkProcedure struct {
	Files  []string
	Into   string
	Target Target
	Linker compilers.Compiler
}

func (p LinkProcedure) RunProcedure() error {
	var err error
	if p.Target == Library {
		err = p.Linker.LinkLibrary(p.Files, p.Into)
	} else {
		err = p.Linker.LinkBinary(p.Files, p.Into)
	}

	if err != nil {
		return err
	}

	// Remove all the object files as we dont really need them.
	for _, s := range p.Files {
		os.Remove(s)
	}

	return nil
}
