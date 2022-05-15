package context

import (
	"chain/compilers"
	"chain/procedures"
	"chain/structures"
	"fmt"
	"path"
)

type Scope struct {
	Parent    *Scope
	Prefix    string
	Libraries []Library
}

func (s *Scope) InheritFrom(parent Scope, prefix string) {
	s.Parent = &parent

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

	compiler := compilers.Clang{
		Path: "/usr/bin/clang",
	}

	buildFiles := []string{}

	for _, f := range procedure.Procedure.Build.Files {
		buildFiles = append(buildFiles, path.Join(s.Prefix, f))
	}

	buildProcedure := procedures.BuildProcedure{
		Files:    buildFiles,
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

	linkProcedure := procedures.LinkProcedure{
		Files:  linkFiles,
		Target: target,
		Into:   path.Join(s.Prefix, procedure.Procedure.Link.Into),
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
}
