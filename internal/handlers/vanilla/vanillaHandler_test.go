package vanilla

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetVersionList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"versions": [
				{"id": "1.16.5", "url": "https://test.com/1.16.5.json"},
				{"id": "1.17", "url": "https://test.com/1.17.json"},
				{"id": "21w19a", "url": "https://test.com/21w19a.json"}
			]
		}`)
	}))

	manifestUrl = mockServer.URL
	defer mockServer.Close()

	versions, err := getVersionsList()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	versionsExpected := Versions{
		Versions: []struct {
			Id  string `json:"id"`
			Url string `json:"url"`
		}{
			{Id: "1.16.5", Url: "https://test.com/1.16.5.json"},
			{Id: "1.17", Url: "https://test.com/1.17.json"},
			{Id: "21w19a", Url: "https://test.com/21w19a.json"},
		},
	}

	if !reflect.DeepEqual(versions, versionsExpected) {
		t.Errorf("Expected versions %v, got %v", versionsExpected, versions)
	}
}

func TestGetUrl(t *testing.T) {
	mux := http.NewServeMux()
	mockServer := httptest.NewServer(mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"versions": [
				{"id": "1.16.5", "url": "%s/1.16.5.json"},
				{"id": "1.17", "url": "%s/1.17.json"},
				{"id": "21w19a", "url": "%s/21w19a.json"}
			]
		}`, mockServer.URL, mockServer.URL, mockServer.URL)
	})

	mux.HandleFunc("/1.16.5.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"downloads": {
				"server": {"url": "http://test.com/1.16.5-server.jar"}
			}
		}`)
	})

	mux.HandleFunc("/1.17.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"downloads": {
				"server": {"url": "https://test.com/1.17-server.jar"}
			}
		}`)
	})

	mux.HandleFunc("/21w19a.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
			"downloads": {
				"server": {"url": ""}
			}
		}`)
	})

	manifestUrl = mockServer.URL
	defer mockServer.Close()

	tests := []struct {
		version     string
		expectedUrl string
		err         bool
	}{
		{"1.16.5", "http://test.com/1.16.5-server.jar", false},
		{"1.17", "https://test.com/1.17-server.jar", false},
		{"21w19a", "", true},
		{"nonexistent", "", true},
	}

	for _, tt := range tests {
		url, err := getUrl(tt.version)
		if tt.err {
			if err == nil {
				t.Errorf("Expected error for version %s, got none", tt.version)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for version %s, got %v", tt.version, err)
			}
			if url != tt.expectedUrl {
				t.Errorf("For version %s, expected URL %s, got %s", tt.version, tt.expectedUrl, url)
			}
		}
	}
}
