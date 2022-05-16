package pkgconfig

import (
	"bytes"
	"os/exec"
	"strings"
)

func CheckPkgExistence(name string) bool {
	args := []string{"--exists", name}

	cmd := exec.Command("pkg-config", args...)

	err := cmd.Run()

	if err != nil {
		return false
	} else {
		return true
	}
}

func GetPkgCflags(name string) ([]string, error) {
	args := []string{"--cflags", name}

	cmd := exec.Command("pkg-config", args...)

	var outbuf bytes.Buffer
	var stdout string

	cmd.Stdout = &outbuf

	err := cmd.Run()

	stdout = outbuf.String()

	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSuffix(stdout, "\n"), " "), nil
}

func GetPkgLibs(name string) ([]string, error) {
	args := []string{"--libs", name}

	cmd := exec.Command("pkg-config", args...)

	var outbuf bytes.Buffer
	var stdout string

	cmd.Stdout = &outbuf

	err := cmd.Run()

	stdout = outbuf.String()

	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSuffix(stdout, "\n"), " "), nil
}

func FindPkg(name string) (*Package, error) {
	var pkg Package
	var err error
	var cflags, libs []string

	if !CheckPkgExistence(name) {
		return nil, &PackageError{
			Exists: false,
		}
	}

	pkg.Name = name

	cflags, err = GetPkgCflags(name)

	if err != nil {
		return nil, err
	}

	pkg.Cflags = cflags

	libs, err = GetPkgLibs(name)

	if err != nil {
		return nil, err
	}

	pkg.Libs = libs

	return &pkg, nil
}
