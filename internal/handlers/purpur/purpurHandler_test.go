package purpur

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
				"1.20.1",
				"1.19.4",
				"1.19.3",
				"1.18.2"
			]
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	versions, err := getVersionsList()
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	expectedVersions := []string{
		"1.18.2",
		"1.19.3",
		"1.19.4",
		"1.20.1",
	}

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("got versions %v, expected %v", versions, expectedVersions)
	}
}

func TestGetBuildList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"builds": {
				"all": [
					"2009",
					"2008",
					"2007",
					"2006"
				]
			}
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	builds, err := getBuildList("1.20.1")
	if err != nil {
		t.Errorf("got error %v, expected no error", err)
	}

	expectedBuilds := []string{
		"2006",
		"2007",
		"2008",
		"2009",
	}

	if !reflect.DeepEqual(builds, expectedBuilds) {
		t.Errorf("got builds %v, expected %v", builds, expectedBuilds)
	}
}

func TestGetUrl(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewServer(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"versions": [
				"1.20.1",
				"1.19.4",
				"1.19.3"
			]
		}`)
	})

	mux.HandleFunc("/1.20.1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"builds": {
				"all": [
					"2009",
					"2008",
					"2007"
				]
			}
		}`)
	})

	mux.HandleFunc("/1.19.4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"builds": {
				"all": [
					"1903",
					"1902",
					"1901"
				]
			}
		}`)
	})

	baseUrl = mockServer.URL
	defer mockServer.Close()

	tests := []struct {
		version     string
		build       string
		expectedUrl string
		err         bool
	}{
		{"1.20.1", "latest", mockServer.URL + "/1.20.1/2009/download", false},
		{"1.20.1", "2008", mockServer.URL + "/1.20.1/2008/download", false},
		{"1.19.4", "latest", mockServer.URL + "/1.19.4/1903/download", false},
		{"nonexistent", "latest", "", true},
		{"1.20.1", "nonexistent", "", true},
	}

	for _, tt := range tests {
		url, err := getUrl(tt.version, tt.build)
		if tt.err {
			if err == nil {
				t.Errorf("got no error, expected error for version %s and build %s", tt.version, tt.build)
			}
		} else {
			if err != nil {
				t.Errorf("got error %v, expected no error for version %s and build %s", err, tt.version, tt.build)
			}
			if url != tt.expectedUrl {
				t.Errorf("got %s, expected %s", url, tt.expectedUrl)
			}
		}
	}
}
