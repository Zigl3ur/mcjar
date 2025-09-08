package paper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetVersionsList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"versions": [
				{
					"version": {
						"id": "1.20.1",
						"support": {
							"status": "supported"
						}
					},
					"builds": [497, 496, 495]
				},
				{
					"version": {
						"id": "1.19.4",
						"support": {
							"status": "supported"
						}
					},
					"builds": [550, 549, 548]
				},
				{
					"version": {
						"id": "23w31a",
						"support": {
							"status": "experimental"
						}
					},
					"builds": [12, 11, 10]
				}
			]
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	versions, err := getVersionsList("paper")
	if err != nil {
		t.Errorf("got error %v, expected no error", err)
	}

	expectedVersions := PaperVersions{
		Versions: []struct {
			Version struct {
				Id      string `json:"id"`
				Support struct {
					Status string `json:"status"`
				} `json:"support"`
			} `json:"version"`
			Builds []int `json:"builds"`
		}{
			{
				Version: struct {
					Id      string `json:"id"`
					Support struct {
						Status string `json:"status"`
					} `json:"support"`
				}{
					Id: "1.20.1",
					Support: struct {
						Status string `json:"status"`
					}{
						Status: "supported",
					},
				},
				Builds: []int{497, 496, 495},
			},
			{
				Version: struct {
					Id      string `json:"id"`
					Support struct {
						Status string `json:"status"`
					} `json:"support"`
				}{
					Id: "1.19.4",
					Support: struct {
						Status string `json:"status"`
					}{
						Status: "supported",
					},
				},
				Builds: []int{550, 549, 548},
			},
			{
				Version: struct {
					Id      string `json:"id"`
					Support struct {
						Status string `json:"status"`
					} `json:"support"`
				}{
					Id: "23w31a",
					Support: struct {
						Status string `json:"status"`
					}{
						Status: "experimental",
					},
				},
				Builds: []int{12, 11, 10},
			},
		},
	}

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("got %v, expected %v", versions, expectedVersions)
	}
}

func TestGetBuildList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/paper/versions/1.20.1/builds" && r.URL.RawQuery == "channel=STABLE" {
			fmt.Fprintln(w, `[
				{"id": 497},
				{"id": 496},
				{"id": 495}
			]`)
		} else if r.URL.Path == "/paper/versions/nonexistent/builds" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{"error": "Version not found"}`)
		}
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	builds, err := getBuildList("paper", "1.20.1")
	if err != nil {
		t.Errorf("got error %v, expected no error", err)
	}

	expectedBuilds := []PaperBuild{
		{Id: 497},
		{Id: 496},
		{Id: 495},
	}

	if !reflect.DeepEqual(builds, expectedBuilds) {
		t.Errorf("got builds %v, expected %v", builds, expectedBuilds)
	}

	if len(builds) != 3 {
		t.Errorf("got %d builds, expected 3", len(builds))
	}

	if builds[0].Id != 497 {
		t.Errorf("got first build ID %d, expected 497", builds[0].Id)
	}

	_, err = getBuildList("paper", "nonexistent")
	if err == nil {
		t.Error("got no error, expected an error")
	}
}

func TestGetUrl(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/paper/versions/1.20.1/builds/latest":
			fmt.Fprintln(w, `{
				"downloads": {
					"server:default": {
						"url": "https://api.papermc.io/v2/projects/paper/versions/1.20.1/builds/497/downloads/paper-1.20.1-497.jar"
					}
				}
			}`)
		case "/paper/versions/1.20.1/builds/496":
			fmt.Fprintln(w, `{
				"downloads": {
					"server:default": {
						"url": "https://api.papermc.io/v2/projects/paper/versions/1.20.1/builds/496/downloads/paper-1.20.1-496.jar"
					}
				}
			}`)
		case "/paper/versions/nonexistent/builds/latest":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{"error": "Version not found"}`)
		case "/paper/versions/1.20.1/builds/999":
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, `{"error": "Build not found"}`)
		}
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	tests := []struct {
		project     string
		version     string
		build       string
		expectedUrl string
		err         bool
	}{
		{"paper", "1.20.1", "latest", "https://api.papermc.io/v2/projects/paper/versions/1.20.1/builds/497/downloads/paper-1.20.1-497.jar", false},
		{"paper", "1.20.1", "496", "https://api.papermc.io/v2/projects/paper/versions/1.20.1/builds/496/downloads/paper-1.20.1-496.jar", false},
		{"paper", "nonexistent", "latest", "", true},
		{"paper", "1.20.1", "999", "", true},
	}

	for _, tt := range tests {
		url, err := getUrl(tt.project, tt.version, tt.build)
		if tt.err {
			if err == nil {
				t.Errorf("got no error, expected error for project %s, version %s and build %s", tt.project, tt.version, tt.build)
			}
		} else {
			if err != nil {
				t.Errorf("got error %v, expected no error for project %s, version %s and build %s", err, tt.project, tt.version, tt.build)
			}
			if url != tt.expectedUrl {
				t.Errorf("got %s, expected %s", url, tt.expectedUrl)
			}
		}
	}
}
