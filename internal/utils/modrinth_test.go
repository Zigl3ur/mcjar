package utils

import "testing"

func TestFacetsBuilder(t *testing.T) {
	type expectedData struct {
		loader      string
		projectType string
		versions    []string
	}

	test := []struct {
		given    expectedData
		expected string
	}{
		{
			expectedData{"forge",
				"mod",
				[]string{"1.20.1", "1.11.2", "1.16.1"}}, "[[\"versions:1.20.1\",\"versions:1.11.2\",\"versions:1.16.1\"],[\"categories:forge\"],[\"project_type:mod\"]]",
		}, {
			expectedData{"",
				"plugin",
				[]string{"1.18.2", "1.8.9", "1.17.1", "1.14.2", "1.7.1"}}, "[[\"versions:1.18.2\",\"versions:1.8.9\",\"versions:1.17.1\",\"versions:1.14.2\",\"versions:1.7.1\"],[\"project_type:plugin\"]]",
		}, {
			expectedData{"fabric",
				"mod",
				[]string{"1.21.1", "1.9.7"}}, "[[\"versions:1.21.1\",\"versions:1.9.7\"],[\"categories:fabric\"],[\"project_type:mod\"]]",
		}, {
			expectedData{"neoforge",
				"mod",
				[]string{"1.9", "1.8.2", "1.15"}}, "[[\"versions:1.9\",\"versions:1.8.2\",\"versions:1.15\"],[\"categories:neoforge\"],[\"project_type:mod\"]]",
		}, {
			expectedData{"forge",
				"modpack",
				[]string{"1.16.3", "1.13.2", "1.8"}}, "[[\"versions:1.16.3\",\"versions:1.13.2\",\"versions:1.8\"],[\"categories:forge\"],[\"project_type:modpack\"]]",
		},
	}

	for _, tt := range test {
		result := FacetsBuilder(tt.given.versions, tt.given.loader, tt.given.projectType)
		if result != tt.expected {
			t.Errorf("got %s, expected %s", result, tt.expected)
		}
	}
}
