package compilers

type Compiler interface {
	Compile(in string, out string) int
}
