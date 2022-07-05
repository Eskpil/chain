package procedures

import (
	"chain/compilers"
	"path"
	"path/filepath"
	"strings"
)

type BuildProcedure struct {
	Files    []string
	Raw      []string
	Cflags   []string
	BuildDir string
	Compiler compilers.Compiler
}

func (p BuildProcedure) RunProcedure() error {
	for i, s := range p.Files {
		output := strings.TrimSuffix(p.Raw[i], filepath.Ext(p.Raw[i])) + ".o"

		output = strings.Replace(output, "/", "_", 12)

		p.Compiler.Compile(s, path.Join(p.BuildDir, output), p.Cflags)
	}

	return nil
}
