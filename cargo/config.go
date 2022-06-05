package cargo

import (
	"chain/logger"
	"fmt"
	"os"
	"os/exec"
)

type CargoConfig struct {
	Name         string
	Path         string
	ConfigPath   string
	Dependencies []string
}

func (config *CargoConfig) GeneratePackageSection(file *os.File) {
	file.WriteString("[package]\n")
	file.WriteString(fmt.Sprintf("name = \"%s\"\n", config.Name))
	file.WriteString(fmt.Sprintf("version = \"0.1.0\"\n"))
	file.WriteString(fmt.Sprintf("edition = \"2021\"\n"))
}

func (config *CargoConfig) GenerateBinSection(file *os.File) {
	file.WriteString("[[bin]]\n")
	file.WriteString(fmt.Sprintf("name = \"%s\"\n", config.Name))
	file.WriteString(fmt.Sprintf("path = \"%s\"\n", config.Path))
}

func (config *CargoConfig) GenerateDependenciesSection(file *os.File) {
	file.WriteString("[dependencies]\n")
	for _, s := range config.Dependencies {
		file.WriteString(fmt.Sprintf("%s = \"*\"\n", s))
	}
}

func (config *CargoConfig) GenerateInto(path string) {
	file, err := os.Create(path)

	if err != nil {
		logger.Error.Printf("Failed to create file: %s\n", path)
		logger.PrintError(fmt.Sprintf("%s", err))
	}

	config.ConfigPath = path

	config.GeneratePackageSection(file)
	config.GenerateBinSection(file)
	config.GenerateDependenciesSection(file)

	file.Close()
}

func (config *CargoConfig) Execute(into string) {
	options := []string{
		"--release",
		"--manifest-path",
		fmt.Sprintf("%s/Cargo.toml", into),
		"--target-dir",
		fmt.Sprintf("%s/garbage", into)}
	env := []string{fmt.Sprintf("CARGO_HOME=%s", fmt.Sprintf("%s/home", into))}

	args := []string{"env"}
	args = append(args, env...)
	args = append(args, "cargo")
	args = append(args, "build")
	args = append(args, options...)

	cmd := exec.Command(args[0], args...)

	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when building cargo project:\n")
		logger.PrintError(string(output))
		os.Exit(1)
	}

	cwd, err := os.Getwd()

	args = []string{
		fmt.Sprintf("%s/%s/garbage/release/%s", cwd, into, config.Name),
		fmt.Sprintf("%s/%s", cwd, into)}

	cmd = exec.Command("cp", args...)

	cmd.Env = os.Environ()

	output, err = cmd.CombinedOutput()

	if err != nil {
		logger.Error.Printf("Error when copying file:\n")
		logger.PrintError(string(output))
		os.Exit(1)
	}
}
