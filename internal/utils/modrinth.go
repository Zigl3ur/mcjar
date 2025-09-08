package utils

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// facets builder create facets string with given data,
// facets is a query param for modrinth api
func FacetsBuilder(versions []string, loader, projectType string) string {
	elt := make([]string, 0, 3)

	if len(versions) > 0 {
		velt := make([]string, 0, len(versions))

		for _, v := range versions {
			velt = append(velt, fmt.Sprintf("\"versions:%s\"", v))
		}

		elt = append(elt, fmt.Sprintf("[%s]", strings.Join(velt, ",")))
	}

	if loader != "" {
		elt = append(elt, fmt.Sprintf("[\"categories:%s\"]", loader))
	}

	if projectType != "" {
		elt = append(elt, fmt.Sprintf("[\"project_type:%s\"]", projectType))
	}

	if len(elt) == 0 {
		return ""
	}

	return fmt.Sprintf("[%s]", strings.Join(elt, ","))
}

type ModsIndex struct {
	Files []struct {
		Env struct {
			Client string `json:"client"`
			Server string `json:"server"`
		} `json:"env"`
		Downloads []string `json:"downloads"`
	} `json:"files"`
}

func MrPackHandler(packPath, modsDir string, isVerbose bool) error {

	uuid := uuid.New()
	output := filepath.Join(os.TempDir(), fmt.Sprintf("mcjar-%s", uuid))

	_ = os.MkdirAll(output, 0755)

	if err := extractIndexJson(packPath, output); err != nil {
		return err
	}

	_, fileMrpack := filepath.Split(packPath)

	//nolint:errcheck
	defer os.Remove(filepath.Join(modsDir, fileMrpack))

	//nolint:errcheck
	defer os.RemoveAll(output)

	modsIndexPath := fmt.Sprintf("%s/modrinth.index.json", output)

	modsIndex, err := os.Open(modsIndexPath)
	if err != nil {
		return errors.New("failed to open modrinth.index.json file")
	}

	var modsData ModsIndex
	if err := json.NewDecoder(modsIndex).Decode(&modsData); err != nil {
		return err
	}

	for _, d := range modsData.Files {
		for _, urlDownload := range d.Downloads {
			urlDownload, _ = url.QueryUnescape(urlDownload)
			parsedUrl := strings.Split(urlDownload, "/")
			filename := parsedUrl[len(parsedUrl)-1]
			if d.Env.Server == "required" {
				if err := WriteToFs(urlDownload, filepath.Join(modsDir, filename)); err != nil {
					fmt.Printf("Failed to get %s", filename)
				}
			}
		}
	}

	return nil
}

func extractIndexJson(mrpackPatch, output string) error {

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
