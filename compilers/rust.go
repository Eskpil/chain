package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Rust struct {
	Path string
}

func (c Rust) Compile(in string, out string, cflags []string) error {

	out = strings.Split(out, ".")[0]

	args := []string{"-o", out, in}

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when compiling file: %s:\n%s", in, output)
		return err
	}

	return nil
}

func (c Rust) LinkBinary(in []string, out string, libraries []Library) error {
	args := []string{"-o", out}

	// TODO: Link with libraries.

	args = append(args, in...)

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when linking: %s\n", output)
		return err
	}

	return nil
}

func (c Rust) LinkLibrary(in []string, out string, libraries []Library) error {
	args := []string{"-shared", "-undefined", "dynamic_lookup", "-o", out}

	for _, library := range libraries {
		for _, flag := range library.Libs {
			args = append(args, flag)
		}
	}

	args = append(args, in...)

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when linking dynamic library: %s\n", output)
		return err
	}

	return nil
}
