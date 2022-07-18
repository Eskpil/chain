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
	Parent *Scope

	Prefix     string
	BuildDir   string
	ObjectsDir string

	Libraries map[string]compilers.Library
	Compilers map[string]structures.Compiler

	Root  bool
	Hooks []Hook
}

func (s *Scope) InheritFrom(parent *Scope, prefix string) {
	cwd, err := os.Getwd()

	if err != nil {
		logger.Error.Println("Failed to get curent working directory.", err)
	}

	s.Prefix = path.Join(parent.Prefix, prefix)
	s.BuildDir = path.Join(parent.BuildDir, prefix)
	s.ObjectsDir = path.Join(parent.BuildDir, prefix, "__objects__")

	if _, err := os.Stat(path.Join(cwd, parent.BuildDir)); os.IsNotExist(err) {
		os.Mkdir(parent.BuildDir, 0777)
	}

	if _, err := os.Stat(path.Join(cwd, s.BuildDir)); os.IsNotExist(err) {
		os.Mkdir(s.BuildDir, 0777)
	}

	if _, err := os.Stat(path.Join(cwd, s.ObjectsDir)); os.IsNotExist(err) {
		os.Mkdir(path.Join(cwd, s.ObjectsDir), 0777)
	}

	s.Compilers = parent.Compilers
	s.Libraries = make(map[string]compilers.Library)
	s.Root = false

	for _, h := range parent.Hooks {
		s.Hooks = append(s.Hooks, h)
	}

	s.Parent = parent
	for _, l := range parent.Libraries {
		s.Libraries[l.Name] = l
	}
}

func (s *Scope) ExportUpwards() {
	for _, l := range s.Libraries {
		s.Parent.Libraries[l.Name] = l
	}
}

func (s *Scope) ExportLibrary(name string) {
	logger.Info.Printf("Exporting: %s\n", name)
	if s.Parent == nil {
		fmt.Println("Current scope does not have a parent.")
		return
	}

	var library *compilers.Library

	for key, lib := range s.Libraries {
		if key == name {
			library = &lib
		}
	}

	if library == nil {
		logger.Error.Printf("Library: %s not found in current scope\n", name)
		os.Exit(1)
	}

	s.Parent.Libraries[library.Name] = *library
}

func (s Scope) FindLibrary(name string) *compilers.Library {
	var library *compilers.Library

	for key, lib := range s.Libraries {
		if key == name {
			library = &lib
		}
	}

	if library == nil {
		logger.Error.Printf("Library: %s not found in current scope\n", name)
		os.Exit(1)
	}

	return library
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
