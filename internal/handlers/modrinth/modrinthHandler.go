package modrinth

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Zigl3ur/mcli/internal/cli/flags"
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
	ClientSide      string `json:"client_side"`
	ServerSide      string `json:"server_side"`
	LoadersVersions map[string][]string
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	CreatedAt       string   `json:"published"`
	UpdatedAt       string   `json:"updated"`
	Downloads       int      `json:"downloads"`
	Categories      []string `json:"categories"`
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

	if status, err := utils.GetReqJson(searchUrl, &results); err != nil {
		return results, fmt.Errorf("failed to fetch search result for %s from Modrinth API (HTTP %d): %w", query, status, err)

	}

	return results, nil
}

func Info(slug string) (SlugData, error) {
	var slugData SlugData

	url := fmt.Sprintf("https://api.modrinth.com/v2/project/%s", slug)

	if status, err := utils.GetReqJson(url, &slugData); err != nil {
		if status == http.StatusNotFound {
			return slugData, fmt.Errorf("no info found for \"%s\" (check slug): %w", slug, err)
		}
		return slugData, fmt.Errorf("failed to fetch %s info from Modrinth API (HTTP %d): %w", slug, status, err)
	}

	var downloadData []DownloadData
	loadersVersions := make(map[string][]string)
	if status, err := utils.GetReqJson(url+"/version", &downloadData); err != nil {
		if status == http.StatusNotFound {
			return slugData, fmt.Errorf("no download data found for \"%s\" (check slug): %w", slug, err)
		}
		return slugData, fmt.Errorf("failed to fetch %s download data from Modrinth API (HTTP %d): %w", slug, status, err)
	}

	for _, data := range downloadData {
		for _, loaderName := range data.Loaders {
			if slices.Contains(flags.ValidLoaders, loaderName) {
				for _, gameVersion := range data.GameVersions {
					if !slices.Contains(loadersVersions[loaderName], gameVersion) {
						loadersVersions[loaderName] = append(loadersVersions[loaderName], gameVersion)
					}
				}
			}
		}
	}

	slugData.LoadersVersions = loadersVersions

	return slugData, nil
}

func Download(slug, version, loader, dir string) (string, error) {
	var downloadData []DownloadData

	if status, err := utils.GetReqJson(fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", slug), &downloadData); err != nil {
		if status == http.StatusNotFound {
			return "", fmt.Errorf("no downloads found for \"%s\" (check slug): %w", slug, err)
		}
		return "", fmt.Errorf("failed to fetch %s downloads from Modrinth API (HTTP %d): %w", slug, status, err)
	}

	idx := -1
	for i, data := range downloadData {
		if slices.Contains(data.Loaders, loader) && slices.Contains(data.GameVersions, version) {
			idx = i
			break
		}
	}

	if idx == -1 {
		return "", fmt.Errorf("no suitable mod version found for %s (loader: %s, game-version: %s)", slug, loader, version)
	}

	var filePath string
	for _, file := range downloadData[idx].Files {
		formattedFilename := strings.ReplaceAll(file.Filename, " ", "_")
		fullPath := filepath.Join(dir, formattedFilename)

		// get filepath for mrpack cause need to extract it
		if filepath.Ext(file.Filename) == ".mrpack" {
			filePath = fullPath
		}
		if err := utils.WriteToFs(file.Url, fullPath); err != nil {
			return "", err
		}
	}

	if len(downloadData[idx].Dependencies) > 0 {
		for _, deps := range downloadData[idx].Dependencies {
			if deps.DependencyType == "required" {
				if _, err := Download(deps.ProjectId, version, loader, dir); err != nil {
					return "", err
				}
			}
		}
	}

	return filePath, nil
}
