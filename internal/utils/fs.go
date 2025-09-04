package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

func ExtractIndexJson(mrpackPatch, output string) error {

	archive, err := zip.OpenReader(mrpackPatch)
	if err != nil {
		return err
	}

	//nolint:errcheck
	defer archive.Close()

	for _, f := range archive.File {
		if f.Name == "modrinth.index.json" {
			filePath := filepath.Join(output, f.Name)

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}

			//nolint:errcheck
			defer dstFile.Close()

			srcFile, err := f.Open()
			if err != nil {
				return err
			}

			//nolint:errcheck
			defer srcFile.Close()

			if _, err := io.Copy(dstFile, srcFile); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("no modrinth.index.json file found")
}
