package pkgconfig

type Package struct {
	Name   string
	Libs   []string
	Cflags []string
}

type PackageError struct {
	Exists bool
}

func (err *PackageError) Error() string {
	if !err.Exists {
		return "Package does not exist."
	}

	return ""
}
