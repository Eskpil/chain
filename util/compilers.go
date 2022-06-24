package util

import "chain/structures"

func LoadDefaultCompilers() structures.CompilersStructure {
	clang := structures.Compiler{
		Name:     "clang",
		Path:     "/usr/bin/clang",
		Language: "c/c++",
		Flags:    []string{"-ggdb", "-O0"},
	}

	linker := structures.Compiler{
		Name:     "clang",
		Path:     "/usr/bin/clang",
		Language: "c/c++/rust",
		Flags:    []string{},
	}

	rustc := structures.Compiler{
		Name:     "rust",
		Path:     "/usr/bin/rustc",
		Language: "rust",
		Flags:    []string{},
	}

	structure := structures.CompilersStructure{}

	structure.Compilers = append(structure.Compilers, clang)
	structure.Compilers = append(structure.Compilers, linker)
	structure.Compilers = append(structure.Compilers, rustc)

	return structure
}
