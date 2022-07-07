package compilers

import (
	"chain/logger"
	"os"
	"os/exec"
)

type Clang struct {
	Path  string
	Flags []string
}

func (c Clang) Language() string {
	return "c/c++"
}

func (c Clang) Compile(in string, out string, cflags []string) error {
	logger.Info.Printf("Compiling: %s\n", in)
	args := []string{"-o", out, "-c", in}

	for _, flag := range c.Flags {
		args = append(args, flag)
	}

	for _, flag := range cflags {
		args = append(args, flag)
	}

	command := exec.Command(c.Path, args...)

	command.Env = os.Environ()

	output, err := command.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when compiling file: %s\n", in)
		logger.PrintError(string(output))
		os.Exit(1)
	}

	return nil
}

func (c Clang) LinkBinary(in []string, out string, libraries []Library) error {
	logger.Info.Printf("Linking binary: %s\n", out)
	args := []string{"-o", out}

	for _, flag := range c.Flags {
		args = append(args, flag)
	}

	for _, library := range libraries {
		for _, flag := range library.Libs {
			args = append(args, flag)
		}
	}

	args = append(args, in...)

	command := exec.Command(c.Path, args...)

	command.Env = os.Environ()

	output, err := command.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when linking binary: %s\n", out)
		logger.Error.Println("Error: ", err)
		logger.PrintError(string(output))
		os.Exit(1)
	}

	return nil
}

func (c Clang) LinkLibrary(in []string, out string, libraries []Library) error {
	logger.Info.Printf("Linking library: %s\n", out)
	args := []string{"-shared", "-o", out}

	for _, flag := range c.Flags {
		args = append(args, flag)
	}

	for _, library := range libraries {
		for _, flag := range library.Libs {
			args = append(args, flag)
		}
	}

	args = append(args, in...)

	command := exec.Command(c.Path, args...)

	command.Env = os.Environ()

	output, err := command.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when linking dynamic library: %s\n", out)
		logger.PrintError(string(output))
		os.Exit(1)
	}

	return nil
}
