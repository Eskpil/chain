package context

import (
	"chain/cargo"
	"chain/compilers"
	"chain/logger"
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
	Compilers map[string]structures.Compiler

	Root  bool
	Hooks []Hook
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

	s.Compilers = parent.Compilers
	s.Root = false

	for _, h := range parent.Hooks {
		s.Hooks = append(s.Hooks, h)
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

func (s Scope) CargoProject(
	libraries []compilers.Library,
	procedure structures.ProcedureStructure) {
	config := cargo.CargoConfig{}

	cwd, err := os.Getwd()

	if err != nil {
		logger.Error.Printf("Unable to get current working directory\n")
		logger.PrintError(fmt.Sprintf("%s", err))
	}

	config.Path = path.Join(s.Prefix, procedure.Procedure.Build.Files[0])
	config.Path = path.Join(cwd, config.Path)
	config.Name = *procedure.Procedure.Name

	for _, lib := range libraries {
		config.Dependencies = append(config.Dependencies, lib.Name)
	}

	config.GenerateInto(fmt.Sprintf("%s/Cargo.toml", s.BuildDir))
	config.Execute(s.BuildDir)
}
