package procedures

import (
	"chain/compilers"
	"fmt"
	"path"
	"strings"
)

type BuildProcedure struct {
	Files    []string
	Cflags   []string
	BuildDir string
	Compiler compilers.Compiler
}

func (p BuildProcedure) RunProcedure() error {
	for _, s := range p.Files {
		raw := strings.Split(path.Base(s), ".")
		output := path.Join(p.BuildDir, fmt.Sprintf("%s.o", raw[0]))

		p.Compiler.Compile(s, output, p.Cflags)
	}

	return nil
}
