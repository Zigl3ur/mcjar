package fabric

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
			"game": [
				{"version": "1.20.1", "stable": true},
				{"version": "1.19.4", "stable": true},
				{"version": "23w31a", "stable": false}
			]
		}`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	versions, err := getVersionsList()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	versionsExpected := FabricVersion{
		Versions: []struct {
			Version string `json:"version"`
			Stable  bool   `json:"stable"`
		}{
			{Version: "1.20.1", Stable: true},
			{Version: "1.19.4", Stable: true},
			{Version: "23w31a", Stable: false},
		},
	}

	if !reflect.DeepEqual(versions, versionsExpected) {
		t.Errorf("got %v, expected %v", versions, versionsExpected)
	}
}

func TestGetStableLoader(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/loader" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `[
				{"version": "0.14.21", "stable": true},
				{"version": "0.14.22", "stable": false},
				{"version": "0.14.20", "stable": true}
			]`)
		}
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	loader, err := getStableLoader()
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	expectedLoader := "0.14.21"
	if loader != expectedLoader {
		t.Errorf("got %s, expected loader %s", loader, expectedLoader)
	}
}

func TestGetStableInstaller(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/versions/installer" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, `[
				{"version": "0.11.2", "stable": true},
				{"version": "0.11.3", "stable": false},
				{"version": "0.11.1", "stable": true}
			]`)
		}
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	installer, err := getStableInstaller()
	if err != nil {
		t.Errorf("got %v, expected no error", err)
	}

	expectedInstaller := "0.11.2"
	if installer != expectedInstaller {
		t.Errorf("got %s, expected installer %s", installer, expectedInstaller)
	}
}

func TestGetUrl(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewServer(mux)

	mux.HandleFunc("/loader", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"version": "0.14.21", "stable": true}]`)
	})

	mux.HandleFunc("/versions/installer", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `[{"version": "0.11.2", "stable": true}]`)
	})

	baseUrl = mockServer.URL
	defer mockServer.Close()

	tests := []struct {
		version     string
		expectedUrl string
	}{
		{"1.20.1", mockServer.URL + "/loader/1.20.1/0.14.21/0.11.2/server/jar"},
		{"1.19.4", mockServer.URL + "/loader/1.19.4/0.14.21/0.11.2/server/jar"},
	}

	for _, tt := range tests {
		url, err := getUrl(tt.version)
		if err != nil {
			t.Errorf("got %v, expected no error for version %s", err, tt.version)
		}
		if url != tt.expectedUrl {
			t.Errorf("got %s, expected URL %s", url, tt.expectedUrl)
		}
	}
}
