package forge

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetVersionsList(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
		<metadata>
			<versioning>
				<latest>1.20.1-47.1.0</latest>
				<release>1.20.1-47.1.0</release>
				<versions>
					<version>1.20.1-47.0.35</version>
					<version>1.20.1-47.1.0</version>
					<version>1.19.4-45.1.0</version>
					<version>1.19.4-45.2.0</version>
				</versions>
			</versioning>
		</metadata>`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	versions, err := getVersionsList()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedVersions := map[string][]string{
		"1.20.1": {"47.1.0", "47.0.35"},
		"1.19.4": {"45.2.0", "45.1.0"},
	}

	if !reflect.DeepEqual(versions, expectedVersions) {
		t.Errorf("got %v, expected %v", versions, expectedVersions)
	}
}

func TestGetUrl(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8"?>
		<metadata>
			<versioning>
				<latest>1.20.1-47.1.0</latest>
				<release>1.20.1-47.1.0</release>
				<versions>
					<version>1.20.1-47.0.35</version>
					<version>1.20.1-47.1.0</version>
					<version>1.19.4-45.1.0</version>
					<version>1.19.4-45.2.0</version>
				</versions>
			</versioning>
		</metadata>`)
	}))

	baseUrl = mockServer.URL
	defer mockServer.Close()

	tests := []struct {
		version     string
		build       string
		expectedUrl string
		err         bool
	}{
		{"1.20.1", "latest", mockServer.URL + "/net/minecraftforge/forge/1.20.1-47.1.0/forge-1.20.1-47.1.0-installer.jar", false},
		{"1.20.1", "47.0.35", mockServer.URL + "/net/minecraftforge/forge/1.20.1-47.0.35/forge-1.20.1-47.0.35-installer.jar", false},
		{"1.19.4", "latest", mockServer.URL + "/net/minecraftforge/forge/1.19.4-45.2.0/forge-1.19.4-45.2.0-installer.jar", false},
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
				t.Errorf("got %v, expected no error for version %s and build %s", err, tt.version, tt.build)
			}
			if url != tt.expectedUrl {
				t.Errorf("got %s, expected %s", url, tt.expectedUrl)
			}
		}
	}
}
