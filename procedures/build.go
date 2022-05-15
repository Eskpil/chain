package procedures

import (
	"chain/compilers"
	"fmt"
	"strings"
)

type BuildProcedure struct {
	Files    []string
	Compiler compilers.Compiler
}

func (p BuildProcedure) RunProcedure() error {
	for _, s := range p.Files {
		raw := strings.Split(s, ".")
		output := fmt.Sprintf("%s.o", raw[0])
		p.Compiler.Compile(s, output)
	}

	return nil
}
