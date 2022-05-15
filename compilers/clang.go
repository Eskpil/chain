package compilers

type Clang struct {
	Path string
}

func (c Clang) Compile(in string, out string) int {
	return 1
}
