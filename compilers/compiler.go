package compilers

import (
	"chain/logger"
	"chain/structures"
	"os"
)

type Compiler interface {
	Compile(in string, out string, cflags []string) error
	LinkLibrary(in []string, out string, libraries []Library) error
	LinkBinary(in []string, out string, libraries []Library) error

	Language() string
}

func CompilerFromName(name string) Compiler {
	if name == "rust" {
		rust := Rust{
			Path: "/usr/bin/rustc",
		}

		return rust
	}

	clang := Clang{
		Path: "/usr/bin/clang",
	}

	return clang

}

func ResolvePathSymlink(path string) string {
	fileInfo, err := os.Stat(path)

	if err != nil {
		logger.Error.Printf("Failed to stat path: %s: %s\n", path, err)
		os.Exit(1)
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		res, err := os.Readlink(path)

		if err != nil {
			logger.Error.Println("Failed to resolve symlink: ", err)
			os.Exit(1)
		}

		return res
	} else {
		return path
	}

}

func CompilerFromStructure(structure structures.Compiler) Compiler {
	if structure.Name == "rust" {
		rust := Rust{
			Path:  ResolvePathSymlink(structure.Path),
			Flags: structure.Flags,
		}

		return rust
	}
	clang := Clang{
		Path:  ResolvePathSymlink(structure.Path),
		Flags: structure.Flags,
	}

	return clang
}
