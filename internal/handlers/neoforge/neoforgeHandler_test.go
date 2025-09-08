package neoforge

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
				"20.4.236",
				"20.4.237",
				"20.2.86",
				"20.2.88"
			]
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	versions, err := getVersionsList()
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	expectedVersions := map[string][]string{
		"1.20.4": {"20.4.237", "20.4.236"},
		"1.20.2": {"20.2.88", "20.2.86"},
	}

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("got versions %v, expected %v", versions, expectedVersions)
	}
}

func TestGetUrl(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"versions": [
				"20.4.236",
				"20.4.237",
				"20.2.88"
			]
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	tests := []struct {
		version     string
		build       string
		expectedUrl string
		err         bool
	}{
		{"1.20.4", "latest", mockServer.URL + "/releases/net/neoforged/neoforge/20.4.237/neoforge-20.4.237-installer.jar", false},
		{"1.20.4", "20.4.236", mockServer.URL + "/releases/net/neoforged/neoforge/20.4.236/neoforge-20.4.236-installer.jar", false},
		{"1.20.2", "latest", mockServer.URL + "/releases/net/neoforged/neoforge/20.2.88/neoforge-20.2.88-installer.jar", false},
		{"nonexistent", "latest", "", true},
		{"1.20.4", "nonexistent", "", true},
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
