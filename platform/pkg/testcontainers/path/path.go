package path

import (
	"os"
	"path/filepath"
	"strings"
)

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for dir != "/" && dir != "" {

		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}

		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if strings.Contains(dir, "microservices-factory") {
				return dir, nil
			}
		}

		dir = filepath.Dir(dir)
	}

	return "", os.ErrNotExist
}
