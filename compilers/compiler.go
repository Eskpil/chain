package compilers

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
