package compilers

type Compiler interface {
	Compile(in string, out string) error
	Link(in []string, out string) error
}
