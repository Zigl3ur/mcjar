package modrinth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
	"github.com/Zigl3ur/mcli/internal/utils/loader"
	"github.com/google/uuid"
)

type SearchResult struct {
	Results []struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"hits"`
	Limit     int `json:"limit"`
	TotalHits int `json:"total_hits"`
}

type SlugData struct {
	ClientSide   string   `json:"client_side"`
	ServerSide   string   `json:"server_side"`
	GameVersions []string `json:"game_versions"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	CreatedAt    string   `json:"published"`
	UpdatedAt    string   `json:"updated"`
	Downloads    int      `json:"downloads"`
	Categories   []string `json:"categories"`
	Loaders      []string `json:"loaders"`
}

type DownloadData struct {
	GameVersions []string `json:"game_versions"`
	Loaders      []string `json:"loaders"`
	Name         string   `json:"name"`
	Files        []struct {
		Url      string `json:"url"`
		Filename string `json:"filename"`
	} `json:"files"`
	Dependencies []struct {
		ProjectId      string `json:"project_id"`
		DependencyType string `json:"dependency_type"`
	} `json:"dependencies"`
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

func Search(query, index, facets string, limit int) (SearchResult, error) {
	var results SearchResult

	query = url.QueryEscape(query)
	index = url.QueryEscape(index)

	var searchUrl string
	if facets != "" {
		facets = url.QueryEscape(facets)
		searchUrl = fmt.Sprintf("https://api.modrinth.com/v2/search?query=%s&limit=%d&index=%s&facets=%s", query, limit, index, facets)
	} else {
		searchUrl = fmt.Sprintf("https://api.modrinth.com/v2/search?query=%s&limit=%d&index=%s", query, limit, index)
	}

	if err := utils.GetReqJson(searchUrl, &results); err != nil {
		return results, errors.New("failed to query modrinth api")
	}

	return results, nil
}

func Info(slug string) (SlugData, error) {
	var data SlugData

	url := fmt.Sprintf("https://api.modrinth.com/v2/project/%s", slug)

	if err := utils.GetReqJson(url, &data); err != nil {
		return data, errors.New("failed to query modrinth api (check slug)")
	}

	return data, nil
}

func Download(slug, version, loader, dir string) error {
	var data []DownloadData

	if err := utils.GetReqJson(fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", slug), &data); err != nil {
		return errors.New("failed to query specified slug")
	}

	idx := -1
	for i, d := range data {
		if slices.Contains(d.Loaders, loader) && slices.Contains(d.GameVersions, version) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("no suitable mod version found (loader: %s, game-version: %s)", loader, version)
	}

	for _, f := range data[idx].Files {
		if err := utils.WriteToFs(f.Url, path.Join(dir, f.Filename)); err != nil {
			return fmt.Errorf("failed to download %s", f.Filename)
		}
	}

	if len(data[idx].Dependencies) > 0 {
		for _, d := range data[idx].Dependencies {
			if d.DependencyType == "required" {
				if err := Download(d.ProjectId, version, loader, dir); err != nil {
					fmt.Printf("failed to download dependency (%s)\n", d.ProjectId)
				}
			}
		}
	}

	return nil
}

func MrPackHandler(packPath, modsDir string, isVerbose bool) error {

	unzip, err := utils.GetPath("unzip")
	if err != nil {
		// if no unzip try with tar ?
		return err
	}

	uuid := uuid.New()
	output := fmt.Sprintf("/tmp/mcli-%s", uuid)

	cmd := exec.Command(unzip, packPath, "-d", output)

	if isVerbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		loader.Start("Extracting modpack archive")
	}

	if err = cmd.Run(); err != nil {
		loader.Stop()
		return errors.New("failed to extract modpack")
	}

	//nolint:errcheck
	defer os.RemoveAll(output)

	modsIndexPath := fmt.Sprintf("%s/modrinth.index.json", output)

	fmt.Println(modsIndexPath)
	modsIndex, err := os.Open(modsIndexPath)
	if err != nil {
		return errors.New("failed to open modpack index file")
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
				if err := utils.WriteToFs(urlDownload, path.Join(modsDir, filename)); err != nil {
					fmt.Printf("Failed to get %s", filename)
				}
			}
		}
	}

	return nil
}
