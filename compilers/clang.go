package compilers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Clang struct {
	Path string
}

func (c Clang) Compile(in string, out string) error {
	cmd := exec.Command(c.Path, "-c", "-o", out, in)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when compiling file: %s:\n%s", in, output)
		return err
	}

	return nil
}

func (c Clang) Link(in []string, out string) error {
	args := []string{"-o", out, strings.Join(in, " ")}

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when linking: %s\n", output)
		return err
	}

	return nil
}
