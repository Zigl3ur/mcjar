package modrinth

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Zigl3ur/mcli/internal/utils"
)

type SearchResult struct {
	Results []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"hits"`
	Limit     int `json:"limit"`
	TotalHits int `json:"total_hits"`
}

func Search(query, index, facets string, limit int) (SearchResult, error) {

	var results SearchResult

	query = url.QueryEscape(query)
	index = url.QueryEscape(index)
	facets = url.QueryEscape(facets)

	url := fmt.Sprintf("https://api.modrinth.com/v2/search?query=%s&limit=%d&index=%s&facets=%s", query, limit, index, facets)

	if err := utils.GetReqJson(url, &results); err != nil {
		return results, errors.New("failed to query modrinth api")
	}

	return results, nil
}
