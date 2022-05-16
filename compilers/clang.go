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

func (c Clang) LinkBinary(in []string, out string, libraries []Library) error {

	args := []string{"-o", out, "-Wl"}
	// clang++ Main.cpp -o foo libchaiscript_stdlib-5.3.1.so -Wl,-rpath,/absolute/path

	for _, library := range libraries {
		if len(library.Path) > 0 {
			args = append(args, fmt.Sprintf("-Wl,-rpath,%s", library.Path+"/"))
			args = append(args, library.Path+"/"+library.Target)
		} else {
			args = append(args, "-l")
			args = append(args, library.Name)
		}
	}

	args = append(args, in...)

	fmt.Println("Args: ", args)

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
		args = append(args, "-l")
		args = append(args, library.Path)
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
