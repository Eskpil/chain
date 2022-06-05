package context

import (
	"chain/compilers"
	"chain/logger"
	"chain/pkgconfig"
	"chain/procedures"
	"chain/structures"
	"fmt"
	"os"
	"path"
)

type Scope struct {
	Parent    *Scope
	Prefix    string
	BuildDir  string
	Libraries []compilers.Library
}

func (s *Scope) InheritFrom(parent *Scope, prefix string) {
	s.Prefix = path.Join(parent.Prefix, prefix)
	s.BuildDir = path.Join(parent.BuildDir, prefix)

	if _, err := os.Stat(parent.BuildDir); os.IsNotExist(err) {
		os.Mkdir(parent.BuildDir, 0777)
	}

	if _, err := os.Stat(s.BuildDir); os.IsNotExist(err) {
		os.Mkdir(s.BuildDir, 0777)
	}

	s.Parent = parent
	for _, l := range parent.Libraries {
		s.Libraries = append(s.Libraries, l)
	}
}

func (s *Scope) ExportLibrary(name string) {
	logger.Info.Printf("Exporting: %s\n", name)
	if s.Parent == nil {
		fmt.Println("Current scope does not have a parent.")
		return
	}

	result := -1

	for i, s := range s.Libraries {
		if s.Name == name {
			result = i
		}
	}

	if 0 > result {
		logger.Error.Printf("Library: %s not found in current scope\n", name)
		os.Exit(1)
	}

	library := s.Libraries[result]

	s.Parent.Libraries = append(s.Parent.Libraries, library)
}

func (s Scope) FindLibrary(name string) {
	logger.Info.Printf("Trying to find library: %s in current scope\n", name)
	result := -1

	for i, s := range s.Libraries {
		if s.Name == name {
			result = i
		}
	}

	if 0 > result {
		logger.Error.Printf("Current scope does not contain library: %s\n", name)
		os.Exit(1)
	}

	logger.Info.Printf("Found library: %s in current scope\n", name)

	library := s.Libraries[result]

	s.Parent.Libraries = append(s.Parent.Libraries, library)
}

func (s Scope) RunProcedure(procedure structures.ProcedureStructure) {
	logger.Info.Printf("Running procedure: %s\n", *procedure.Procedure.Name)
	var err error

	libraries := []compilers.Library{}

	if procedure.Procedure.Link != nil {
		for _, with := range procedure.Procedure.Link.With {
			var result *compilers.Library = &compilers.Library{}

			if with.Kind == "exported" {
				for _, library := range s.Libraries {
					if library.Name == with.Name {
						result = &library
					}
				}
			} else if with.Kind == "compiler" {
				result.Name = with.Name
				result.Libs = []string{fmt.Sprintf("-l%s", with.Name)}
			} else if with.Kind == "pkg-config" {
				var pkg *pkgconfig.Package
				pkg, err = pkgconfig.FindPkg(with.Name)

				if err != nil {
					fmt.Printf("Package: %s not found by pkg-config.\n", with.Name)
					return
				}

				result.Libs = pkg.Libs
				result.Name = pkg.Name
				result.Cflags = pkg.Cflags
			} else {
				logger.Error.Printf("Unknown linking method: %s\n", with.Kind)
				os.Exit(1)
			}

			if result != nil {
				libraries = append(libraries, *result)
			} else {
				logger.Error.Printf("Library: %s not found in current scope, have you forgotten to export it?\n", with)
				os.Exit(1)
			}
		}
	}

	compilerName := procedure.Procedure.Build.Compiler
	compiler := compilers.CompilerFromName(compilerName)

	buildFiles := []string{}

	var cflags []string

	for _, library := range libraries {
		for _, flag := range library.Cflags {
			cflags = append(cflags, flag)
		}
	}

	for _, f := range procedure.Procedure.Build.Files {
		buildFiles = append(buildFiles, path.Join(s.Prefix, f))
	}

	buildProcedure := procedures.BuildProcedure{
		Files:    buildFiles,
		BuildDir: s.BuildDir,
		Cflags:   cflags,
		Compiler: compiler,
	}

	err = buildProcedure.RunProcedure()

	if err != nil {
		return
	}

	if procedure.Procedure.Link != nil {
		var target procedures.Target

		if procedure.Procedure.Link.Target == "library" {
			target = procedures.Library
		} else {
			target = procedures.Binary
		}

		linkFiles := []string{}

		for _, f := range procedure.Procedure.Link.Files {
			linkFiles = append(linkFiles, path.Join(s.BuildDir, f))
		}

		linker := compilers.CompilerFromName(procedure.Procedure.Link.Linker)

		libraryPath := path.Join(s.BuildDir, procedure.Procedure.Link.Into)

		linkProcedure := procedures.LinkProcedure{
			Files:  linkFiles,
			Target: target,
			With:   libraries,
			Into:   libraryPath,
			Linker: linker,
		}

		err = linkProcedure.RunProcedure()

		if err != nil {
			return
		}

		if procedure.Procedure.Library != nil && linkProcedure.Target == procedures.Library {
			library := compilers.Library{}

			library.Name = procedure.Procedure.Library.Name

			libpath := s.BuildDir

			library.Libs = []string{fmt.Sprintf("-Wl,-rpath,%s", libpath), fmt.Sprintf("%s/%s.so", libpath, procedure.Procedure.Link.Target)}

			s.Libraries = append(s.Libraries, library)
		}

	}
	if procedure.Procedure.Export != nil {
		for _, e := range procedure.Procedure.Export {
			s.ExportLibrary(e)
		}
	}
}
