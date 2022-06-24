package compilers

import "chain/structures"

type Compiler interface {
	Compile(in string, out string, cflags []string) error
	LinkLibrary(in []string, out string, libraries []Library) error
	LinkBinary(in []string, out string, libraries []Library) error
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

func CompilerFromStructure(structure structures.Compiler) Compiler {
	if structure.Name == "rust" {
		rust := Rust{
			Path:  structure.Path,
			Flags: structure.Flags,
		}

		return rust
	}

	clang := Clang{
		Path:  structure.Path,
		Flags: structure.Flags,
	}

	return clang
}
