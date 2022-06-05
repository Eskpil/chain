package compilers

import (
	"chain/logger"
	"os"
	"os/exec"
)

type Clang struct {
	Path string
}

func (c Clang) Compile(in string, out string, cflags []string) error {
	logger.Info.Printf("Compiling: %s into: %s\n", in, out)
	args := []string{"-o", out, "-c", in}

	for _, flag := range cflags {
		args = append(args, flag)
	}

	cmd := exec.Command(c.Path, args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when compiling file: %s into: %s\n", in, out)
		logger.PrintError(string(output))
		return err
	}

	return nil
}

func (c Clang) LinkBinary(in []string, out string, libraries []Library) error {
	logger.Info.Printf("Linking binary: %s from: %v\n", out, in)
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
		logger.Error.Printf("Error when linking binary: %s\n", out)
		logger.PrintError(string(output))
		return err
	}

	return nil
}

func (c Clang) LinkLibrary(in []string, out string, libraries []Library) error {
	logger.Info.Printf("Linking library: %s from: %v\n", out, in)
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
		logger.Error.Printf("Error when linking dynamic library: %s\n", out)
		logger.PrintError(string(output))
		return err
	}

	return nil
}
