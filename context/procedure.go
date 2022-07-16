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
	"path/filepath"
	"strings"
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

type hookConfig struct {
	Sources string
	Headers string
}

func (s Scope) RunBuildHooks(procedure structures.ProcedureStructure) hookConfig {
	cwd, err := os.Getwd()

	if err != nil {
		logger.Error.Printf("Could not get current working directory: %", err)
		os.Exit(1)
	}

	headersDir := path.Join(cwd, s.BuildDir, "__headers__")
	sourcesDir := path.Join(cwd, s.BuildDir, "__sources__")

	if _, err := os.Stat(headersDir); os.IsNotExist(err) {
		os.Mkdir(headersDir, 0777)
	}

	if _, err := os.Stat(sourcesDir); os.IsNotExist(err) {
		os.Mkdir(sourcesDir, 0777)
	}

	env := os.Environ()

	env = append(env, fmt.Sprintf("CHAIN_HEADERS=%s", headersDir))
	env = append(env, fmt.Sprintf("CHAIN_SOURCES=%s", sourcesDir))

	for _, name := range procedure.Procedure.Build.Hook {
		var hook *Hook

		hook = nil

		for i := range s.Hooks {
			if s.Hooks[i].Name == name {
				hook = &s.Hooks[i]
			}
		}

		if hook == nil {
			logger.Error.Printf("Hook: %s is not found in scope.\n", name)
			os.Exit(1)
		}

		logger.Info.Println("Running build hook: ", hook.Name)

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

	filepath.Walk(sourcesDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error.Println("Failed walk __sources__ directory: ", err)
			os.Exit(1)
		}

		if sourcesDir == filePath {
			return nil
		}

		output := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())) + ".o"

		output = strings.Replace("__"+output, "/", "_", 12)

		compiler := s.DecideCompiler(procedure)

		compiler.Compile(filePath, path.Join(s.ObjectsDir, output), []string{})

		return nil
	})

	config := hookConfig{
		Sources: sourcesDir,
		Headers: headersDir,
	}

	return config
}

func (s Scope) DecideCompiler(procedure structures.ProcedureStructure) compilers.Compiler {
	structure := structures.Compiler{}

	compilerName := procedure.Procedure.Build.Compiler

	for _, comp := range s.Compilers {
		if comp.Name == compilerName {
			structure = comp
		}
	}

	if structure.Path == "" {
		logger.Error.Println("Unable to find compiler: ", compilerName)
		os.Exit(1)
	}

	compiler := compilers.CompilerFromStructure(structure)

	return compiler
}

func (s *Scope) RunBuildSubProcedure(procedure structures.ProcedureStructure) []string {
	config := s.RunBuildHooks(procedure)

	var err error

	libraries := s.GetLibraries(procedure)

	compiler := s.DecideCompiler(procedure)

	var cflags []string

	for _, library := range libraries {
		for _, flag := range library.Cflags {
			cflags = append(cflags, flag)
		}
	}

	if (compiler.Language() == "c/c++") && procedure.Procedure.Build.Headers != nil {
		if *procedure.Procedure.Build.Headers == "." {
			cflags = append(cflags, fmt.Sprintf("-I%s", s.Prefix))
			cflags = append(cflags, fmt.Sprintf("-I%s", config.Headers))
		}
	}

	buildFiles := []string{}
	rawFiles := []string{}

	for _, f := range procedure.Procedure.Build.Files {
		buildFiles = append(buildFiles, path.Join(s.Prefix, f))
		rawFiles = append(rawFiles, f)
	}

	buildDir := s.ObjectsDir

	if compiler.Language() == "rust" {
		buildDir = s.BuildDir
	}

	buildProcedure := procedures.BuildProcedure{
		Files:    buildFiles,
		Raw:      rawFiles,
		BuildDir: buildDir,
		Cflags:   cflags,
		Compiler: compiler,
	}

	err = buildProcedure.RunProcedure()

	if err != nil {
		logger.Error.Printf("Failed to run subprocedure build: ", err)
		os.Exit(1)
	}

	return cflags
}

func (s *Scope) RunLinkSubProcedure(procedure structures.ProcedureStructure, cflags []string) {
	var err error

	libraries := s.GetLibraries(procedure)

	var target procedures.Target

	if procedure.Procedure.Link.Target == "library" {
		target = procedures.Library
	} else {
		target = procedures.Binary
	}

	linkFiles := []string{}

	filepath.Walk(s.ObjectsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error.Println("Failed walk __sources__ directory: ", err)
			os.Exit(1)
		}

		if s.ObjectsDir == path {
			return nil
		}

		linkFiles = append(linkFiles, path)

		return nil
	})

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

		library.Cflags = cflags

		libpath := s.BuildDir

		cwd, err := os.Getwd()

		if err != nil {
			logger.Error.Println("Failed to get cwd: ", err)
			os.Exit(1)
		}

		library.Libs = []string{fmt.Sprintf("-Wl,-rpath,%s", libpath), fmt.Sprintf("%s/%s", path.Join(cwd, libpath), procedure.Procedure.Link.Into)}

		s.Libraries = append(s.Libraries, library)
	}

}

func (s *Scope) RunProcedure(procedure structures.ProcedureStructure) {
	logger.Info.Printf("Running procedure: %s\n", *procedure.Procedure.Name)

	cflags := s.RunBuildSubProcedure(procedure)

	if procedure.Procedure.Link != nil {
		s.RunLinkSubProcedure(procedure, cflags)
	}

	if procedure.Procedure.Export != nil {
		for _, e := range procedure.Procedure.Export {
			s.ExportLibrary(e)
		}
	}
}
