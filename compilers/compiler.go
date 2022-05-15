package compilers

type Compiler interface {
	Compile(in string, out string) error
	LinkLibrary(in []string, out string, libraries []Library) error
	LinkBinary(in []string, out string, libraries []Library) error
}
