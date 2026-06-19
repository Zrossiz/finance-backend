package helpers

import (
	"fmt"
	"os"
)

func WithBasePath(path string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", pwd, path), nil
}
