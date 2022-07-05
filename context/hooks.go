package context

import (
	"chain/logger"
	"fmt"
	"os"
	"path/filepath"
)

type Hook struct {
	Name string
	Path string
}

func (project *Scope) LoadHooks() {
	cwd, err := os.Getwd()

	if err != nil {
		logger.Error.Println("Failed to get current working directory: ", cwd)
		os.Exit(1)
	}

	hooksDir := fmt.Sprintf("%s/.hooks", cwd)

	fileInfo, err := os.Stat(hooksDir)

	if err != nil {
		if os.IsNotExist(err) {
			// Dosen't matter. User might not have defined any hooks.
			return
		}

		logger.Error.Println("Should handle error: ", err)
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		logger.Error.Printf("Expected: %s to be a directory.\n", hooksDir)
		os.Exit(1)
	}

	filepath.Walk(hooksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error.Println("Failed walk hooks directory: ", err)
			os.Exit(1)
		}

		hook := Hook{}

		hook.Name = info.Name()
		hook.Path = path

		project.Hooks = append(project.Hooks, hook)

		return nil
	})
}
