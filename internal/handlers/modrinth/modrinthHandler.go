package modrinth

import (
	"errors"
	"fmt"
	"net/url"

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
		Url string `json:"url"`
	} `json:"files"`
	Dependencies []struct {
		ProjectId      string `json:"project_id"`
		DependencyType string `json:"required"`
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

func Download(slug, version, loader, path string) error {
	var data []DownloadData

	if err := utils.GetReqJson(fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", slug), &data); err != nil {
		return errors.New("failed to query specified slug")
	}

	if loader == "" {
		return errors.New("please specify a loader")
	}

	fmt.Println(data)

	// recall to download deps

	return nil
}
