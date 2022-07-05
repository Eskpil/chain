package context

import (
	"chain/compilers"
	"chain/logger"
	"chain/pkgconfig"
	"chain/procedures"
	"chain/structures"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func (s Scope) GetLibraries(procedure structures.ProcedureStructure) []compilers.Library {
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
			} else if with.Kind == "cargo" {
				result.Name = with.Name
			} else if with.Kind == "pkg-config" {
				var pkg *pkgconfig.Package
				pkg, err := pkgconfig.FindPkg(with.Name)

				if err != nil {
					logger.Error.Printf("Package: %s not found by pkg-confg.\n", with.Name)
					os.Exit(1)
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

	return libraries
}

func (s Scope) RunBuildHooks(procedure structures.ProcedureStructure) {
	cwd, err := os.Getwd()

	if err != nil {
		logger.Error.Printf("Could not get current working directory: %", err)
		os.Exit(1)
	}

	headersDir := fmt.Sprintf("%s/%s/__headers__", cwd, s.BuildDir)
	sourcesDir := fmt.Sprintf("%s/%s/__sources__", cwd, s.BuildDir)

	env := os.Environ()

	env = append(env, fmt.Sprintf("CHAIN_HEADERS=%s", headersDir))
	env = append(env, fmt.Sprintf("CHAIN_SOURCES=%s", sourcesDir))

	for _, name := range procedure.Procedure.Build.Hook {
		var hook *Hook

		hook = nil

		for _, x := range s.Hooks {
			if x.Name == name {
				hook = &x
			}
		}

		if hook == nil {
			logger.Error.Printf("Hook: %s is not found in scope.\n", name)
			os.Exit(1)
		}

		logger.Info.Println("Running build hook: ", name)

		args := []string{}
		command := exec.Command(hook.Path, args...)

		command.Env = env

		output, err := command.CombinedOutput()

		if err != nil {
			logger.Error.Printf("Error when executing: %s\n", hook.Name)
			logger.PrintError(string(output))
			os.Exit(1)
		}
	}
}

func (s *Scope) RunBuildSubProcedure(procedure structures.ProcedureStructure) {
	s.RunBuildHooks(procedure)

	var err error

	libraries := s.GetLibraries(procedure)

	structure := structures.Compiler{}

	compilerName := procedure.Procedure.Build.Compiler

	for _, comp := range s.Compilers {
		if comp.Name == compilerName {
			structure = comp
		}
	}

	compiler := compilers.CompilerFromStructure(structure)

	buildFiles := []string{}

	var cflags []string

	for _, library := range libraries {
		for _, flag := range library.Cflags {
			cflags = append(cflags, flag)
		}
	}

	if (structure.Language == "c/c++") && procedure.Procedure.Build.Headers != nil {
		if *procedure.Procedure.Build.Headers == "." {
			cflags = append(cflags, fmt.Sprintf("-I%s", s.Prefix))
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
		logger.Error.Printf("Failed to run subprocedure build: ", err)
		os.Exit(1)
	}
}

func (s *Scope) RunLinkSubProcedure(procedure structures.ProcedureStructure) {
	var err error

	libraries := s.GetLibraries(procedure)

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

	structure := structures.Compiler{}

	linkerName := procedure.Procedure.Link.Linker

	for _, comp := range s.Compilers {
		if comp.Name == linkerName {
			structure = comp
		}
	}

	linker := compilers.CompilerFromStructure(structure)

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

		library.Cflags = []string{fmt.Sprintf("-I%s", s.Prefix)}

		libpath := s.BuildDir

		library.Libs = []string{fmt.Sprintf("-Wl,-rpath,%s", libpath), fmt.Sprintf("%s/%s.so", libpath, procedure.Procedure.Link.Target)}

		s.Libraries = append(s.Libraries, library)
	}

}

func (s *Scope) RunProcedure(procedure structures.ProcedureStructure) {
	logger.Info.Printf("Running procedure: %s\n", *procedure.Procedure.Name)

	s.RunBuildSubProcedure(procedure)

	if procedure.Procedure.Link != nil {
		s.RunLinkSubProcedure(procedure)
	}

	if procedure.Procedure.Export != nil {
		for _, e := range procedure.Procedure.Export {
			s.ExportLibrary(e)
		}
	}
}
