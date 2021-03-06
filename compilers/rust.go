package compilers

import (
	"chain/logger"
	"os"
	"os/exec"
	"strings"
)

type Rust struct {
	Path  string
	Flags []string
}

func (c Rust) Language() string {
	return "rust"
}

func (c Rust) Compile(in string, out string, cflags []string) error {
	out = strings.Split(out, ".")[0]
	args := []string{"-o", out, in}

	for _, flag := range c.Flags {
		args = append(args, flag)
	}

	cmd := exec.Command(c.Path, args...)
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when compiling file: %s\n", in)
		logger.PrintError(string(output))
		logger.Error.Printf("Err: %v\n", err)
		os.Exit(1)
	}

	return nil
}

func (c Rust) LinkBinary(in []string, out string, libraries []Library) error {
	logger.Warn.Println("Linking of binaries not implemented")
	return nil
}

func (c Rust) LinkLibrary(in []string, out string, libraries []Library) error {
	logger.Warn.Printf("Linking of libraries not implemented")
	return nil
}
