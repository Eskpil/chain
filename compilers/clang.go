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

func (c Clang) Compile(in string, out string, cflags []string) error {
	args := []string{"-o", out, "-c", in}

	for _, flag := range cflags {
		args = append(args, flag)
	}

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when compiling file: %s:\n%s", in, output)
		return err
	}

	return nil
}

func (c Clang) LinkBinary(in []string, out string, libraries []Library) error {
	args := []string{"-o", out}

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
		fmt.Printf("Error when linking: %s\n", output)
		return err
	}

	return nil
}

func (c Clang) LinkLibrary(in []string, out string, libraries []Library) error {
	args := []string{"-shared", "-undefined", "dynamic_lookup", "-o", out, strings.Join(in, " ")}

	for _, library := range libraries {
		for _, flag := range library.Cflags {
			args = append(args, flag)
		}
		for _, flag := range library.Cflags {
			args = append(args, flag)
		}
	}

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error when linking dynamic library: %s\n", output)
		return err
	}

	return nil
}
