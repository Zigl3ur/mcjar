package modrinth

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcli/internal/utils"
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

func Download(slug, version, loader, dir string) (string, error) {
	var data []DownloadData

	if err := utils.GetReqJson(fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", slug), &data); err != nil {
		return "", errors.New("failed to query specified slug")
	}

	idx := -1
	for i, d := range data {
		if slices.Contains(d.Loaders, loader) && slices.Contains(d.GameVersions, version) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return "", fmt.Errorf("no suitable mod version found (loader: %s, game-version: %s)", loader, version)
	}

	var filePath string
	for _, f := range data[idx].Files {
		formattedFilename := strings.ReplaceAll(f.Filename, " ", "_")
		fullPath := filepath.Join(dir, formattedFilename)

		// get filepath for mrpack cause need to extract it
		if filepath.Ext(f.Filename) == ".mrpack" {
			filePath = fullPath
		}
		if err := utils.WriteToFs(f.Url, fullPath); err != nil {
			return "", fmt.Errorf("failed to download %s", formattedFilename)
		}
	}

	if len(data[idx].Dependencies) > 0 {
		for _, d := range data[idx].Dependencies {
			if d.DependencyType == "required" {
				if _, err := Download(d.ProjectId, version, loader, dir); err != nil {
					fmt.Printf("failed to download dependency (%s)\n", d.ProjectId)
				}
			}
		}
	}

	return filePath, nil
}
