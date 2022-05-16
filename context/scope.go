package context

import (
	"chain/compilers"
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
	Libraries []compilers.Library
}

func (s *Scope) InheritFrom(parent *Scope, prefix string) {
	s.Parent = parent
	s.Prefix = path.Join(parent.Prefix, prefix)

	for _, l := range parent.Libraries {
		s.Libraries = append(s.Libraries, l)
	}
}

func (s *Scope) ExportLibrary(name string) {
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
		fmt.Println("Current scope does not contain library: ", name)
		return
	}

	library := s.Libraries[result]

	s.Parent.Libraries = append(s.Parent.Libraries, library)
	return
}

func (s Scope) RunProcedure(procedure structures.ProcedureStructure) {
	var err error

	libraries := []compilers.Library{}

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
				return
			}

			result.Libs = pkg.Libs
			result.Name = pkg.Name
			result.Cflags = pkg.Cflags
		}

		if result != nil {
			libraries = append(libraries, *result)
		} else {
			fmt.Printf("Library: %s not found in current scope, have you forgotten to export it?\n", with)
		}
	}

	compiler := compilers.Clang{
		Path: "/usr/bin/clang",
	}

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
		Cflags:   cflags,
		Compiler: compiler,
	}

	var target procedures.Target

	if procedure.Procedure.Link.Target == "library" {
		target = procedures.Library
	} else {
		target = procedures.Binary
	}

	linkFiles := []string{}

	for _, f := range procedure.Procedure.Link.Files {
		linkFiles = append(linkFiles, path.Join(s.Prefix, f))
	}

	libraryPath := path.Join(s.Prefix, procedure.Procedure.Link.Into)

	linkProcedure := procedures.LinkProcedure{
		Files:  linkFiles,
		Target: target,
		With:   libraries,
		Into:   libraryPath,
		Linker: compiler,
	}

	err = buildProcedure.RunProcedure()

	if err != nil {
		return
	}

	err = linkProcedure.RunProcedure()

	if err != nil {
		return
	}

	if procedure.Procedure.Library != nil && linkProcedure.Target == procedures.Library {
		library := compilers.Library{}

		var cwd string

		cwd, err = os.Getwd()

		if err != nil {
			fmt.Println("Failed to get current working directory: ", err)
			return
		}

		library.Name = procedure.Procedure.Library.Name

		libpath := path.Join(cwd, s.Prefix)

		library.Libs = []string{fmt.Sprintf("-Wl,-rpath,%s", libpath), fmt.Sprintf("%s/%s.so", libpath, procedure.Procedure.Link.Target)}

		s.Libraries = append(s.Libraries, library)
	}

	if procedure.Procedure.Export != nil {
		for _, e := range procedure.Procedure.Export {
			s.ExportLibrary(e)
		}
	}
}
