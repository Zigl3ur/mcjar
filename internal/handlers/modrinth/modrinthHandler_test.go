package modrinth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"hits": [
				{
					"slug": "sodium",
					"title": "Sodium",
					"description": "A modern rendering engine for Minecraft"
				},
				{
					"slug": "lithium",
					"title": "Lithium",
					"description": "A general-purpose optimization mod"
				}
			],
			"limit": 10,
			"total_hits": 2
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	result, err := Search("optimization", "relevance", "", 10)
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	expectedResult := SearchResult{
		Results: []struct {
			Slug        string `json:"slug"`
			Title       string `json:"title"`
			Description string `json:"description"`
		}{
			{
				Slug:        "sodium",
				Title:       "Sodium",
				Description: "A modern rendering engine for Minecraft",
			},
			{
				Slug:        "lithium",
				Title:       "Lithium",
				Description: "A general-purpose optimization mod",
			},
		},
		Limit:     10,
		TotalHits: 2,
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("got %v, expected %v", result, expectedResult)
	}
}

func TestSearchWithFacets(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"hits": [
				{
					"slug": "sodium",
					"title": "Sodium",
					"description": "A modern rendering engine for Minecraft"
				}
			],
			"limit": 5,
			"total_hits": 1
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	result, err := Search("rendering", "relevance", "[\"categories:optimization\"]", 5)
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	if len(result.Results) != 1 {
		t.Errorf("got %d, expected 1 result", len(result.Results))
	}

	if result.Results[0].Slug != "sodium" {
		t.Errorf("got '%s', expected slug 'sodium'", result.Results[0].Slug)
	}
}

func TestInfo(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewServer(mux)

	mux.HandleFunc("/project/sodium", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"title": "Sodium",
			"description": "A modern rendering engine for Minecraft",
			"client_side": "required",
			"server_side": "optional",
			"published": "2020-04-27T19:42:33.887573Z",
			"updated": "2023-08-15T10:30:45.123456Z",
			"downloads": 15000000,
			"categories": ["optimization", "utility"]
		}`)
	})

	mux.HandleFunc("/project/sodium/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[
			{
				"game_versions": ["1.20.1", "1.20"],
				"loaders": ["fabric", "forge"],
				"name": "Sodium 0.5.3",
				"files": [
					{
						"url": "https://test.com/sodium-0.5.3.jar",
						"filename": "sodium-0.5.3.jar"
					}
				],
				"dependencies": []
			},
			{
				"game_versions": ["1.19.4", "1.19.3"],
				"loaders": ["fabric"],
				"name": "Sodium 0.4.10",
				"files": [
					{
						"url": "https://test.com/sodium-0.4.10.jar",
						"filename": "sodium-0.4.10.jar"
					}
				],
				"dependencies": []
			},
			{
				"game_versions": ["1.18.2"],
				"loaders": ["quilt"],
				"name": "Sodium 0.4.2",
				"files": [
					{
						"url": "https://test.com/sodium-0.4.2.jar",
						"filename": "sodium-0.4.2.jar"
					}
				],
				"dependencies": []
			}
		]`)
	})

	baseUrl = mockServer.URL
	defer mockServer.Close()

	info, err := Info("sodium")
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	if info.Title != "Sodium" {
		t.Errorf("got '%s', expected title 'Sodium'", info.Title)
	}

	if info.ClientSide != "required" {
		t.Errorf("got '%s', expected client_side 'required'", info.ClientSide)
	}

	if info.ServerSide != "optional" {
		t.Errorf("got '%s', expected server_side 'optional'", info.ServerSide)
	}

	if info.Downloads != 15000000 {
		t.Errorf("got %d, expected downloads 15000000", info.Downloads)
	}

	expectedCategories := []string{"optimization", "utility"}
	if !reflect.DeepEqual(info.Categories, expectedCategories) {
		t.Errorf("got %v, expected categories %v", info.Categories, expectedCategories)
	}

	loadersVersionsExpected := map[string][]string{
		"fabric": {"1.20.1", "1.20", "1.19.4", "1.19.3"},
		"forge":  {"1.20.1", "1.20"},
	}

	if !reflect.DeepEqual(info.LoadersVersions, loadersVersionsExpected) {
		t.Errorf("got %v, expected LoadersVersions %v", info.LoadersVersions, loadersVersionsExpected)
	}
}

func TestInfoNotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, `{"error": "Not Found"}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	_, err := Info("nonexistent-mod")
	if err == nil {
		t.Error("got no error, expected error for non-existent mod")
	}
}
