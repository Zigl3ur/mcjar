package utils

import (
	"errors"
	"fmt"
	"os/exec"
)

func GetPath(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return "", fmt.Errorf("%s not found in PATH, please install it and retry", file)
		} else {
			return "", err
		}
	}

	return path, nil
}
