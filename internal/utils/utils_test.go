package utils

import (
	"slices"
	"testing"
)

type testData[T any, E any] struct {
	given    T
	expected E
	err      bool
}

func TestMcVersionParser(t *testing.T) {
	tests := []struct {
		given    string
		expected [3]int
		unparsed string
	}{
		{"1.8.9", [3]int{1, 8, 9}, ""},
		{"1.12.2", [3]int{1, 12, 2}, ""},
		{"1.10", [3]int{1, 10, 0}, ""},
		{"24w5a", [3]int{0, 0, 0}, "24w5a"},
		{"24w5a.12.ZAE", [3]int{0, 0, 0}, "24w5a.12.ZAE"},
		{"", [3]int{0, 0, 0}, ""},
	}

	for _, tt := range tests {
		result, unparsed := mcVersionParser(tt.given)

		if unparsed != tt.unparsed {
			t.Errorf("got %s, expected %s", unparsed, tt.unparsed)
		}

		if result != tt.expected {
			t.Errorf("got %d, expected %d", result, tt.expected)
		}
	}

}

func TestSortMcVersions(t *testing.T) {
	tests := []testData[[]string, map[string][]string]{
		{[]string{"1.21.6", "1.8.9", "23w6a", "1.12.2", "1.7.2_pre4", "23w8b"}, map[string][]string{"versions": {"1.21.6", "1.12.2", "1.8.9"}, "snapshots": {"23w6a", "1.7.2_pre4", "23w8b"}}, false},
		{[]string{"1.9", "1.10.2", "1.4.3"}, map[string][]string{"versions": {"1.10.2", "1.9", "1.4.3"}, "snapshots": {}}, false},
		{[]string{"1.21.6", "23w6a", "1.21.6"}, map[string][]string{"versions": {"1.21.6", "1.21.6"}, "snapshots": {"23w6a"}}, false},
		{[]string{"21E", " "}, map[string][]string{"versions": {}, "snapshots": {"21E", " "}}, false},
		{[]string{""}, map[string][]string{"versions": {""}, "snapshots": {}}, false},
	}

	for _, tt := range tests {
		result := SortMcVersions(tt.given)
		if !slices.Equal(result["versions"], tt.expected["versions"]) || !slices.Equal(result["snapshots"], tt.expected["snapshots"]) {
			t.Errorf("got %s expected %s", result, tt.expected)
		}
	}

}

func TestHumanizeByte(t *testing.T) {
	tests := []testData[int64, string]{
		{1, "1 B", false},
		{1024, "1.0 KiB", false},
		{1048576, "1.0 MiB", false},
		{1073741824, "1.0 GiB", false},
		{1099511627776, "1.0 TiB", false},
		{1125899906842624, "1.0 PiB", false},
		{512, "512 B", false},
		{000000, "0 B", false},
		{2097152, "2.0 MiB", false},
		{1610612736, "1.5 GiB", false},
		{999, "999 B", false},
	}

	for _, tt := range tests {
		result := humanizeByte(tt.given)
		if result != tt.expected {
			t.Errorf("got %s expected %s", result, tt.expected)
		}
	}
}

func TestFacetsBuilder(t *testing.T) {
	type expectedData struct {
		loader      string
		projectType string
		versions    []string
	}

	test := []testData[expectedData, string]{
		{
			expectedData{"forge",
				"mod",
				[]string{"1.20.1", "1.11.2", "1.16.1"}}, "[[\"versions:1.20.1\",\"versions:1.11.2\",\"versions:1.16.1\"],[\"categories:forge\"],[\"project_type:mod\"]]", false,
		}, {
			expectedData{"",
				"plugin",
				[]string{"1.18.2", "1.8.9", "1.17.1", "1.14.2", "1.7.1"}}, "[[\"versions:1.18.2\",\"versions:1.8.9\",\"versions:1.17.1\",\"versions:1.14.2\",\"versions:1.7.1\"],[\"project_type:plugin\"]]", false,
		}, {
			expectedData{"fabric",
				"mod",
				[]string{"1.21.1", "1.9.7"}}, "[[\"versions:1.21.1\",\"versions:1.9.7\"],[\"categories:fabric\"],[\"project_type:mod\"]]", false,
		}, {
			expectedData{"neoforge",
				"mod",
				[]string{"1.9", "1.8.2", "1.15"}}, "[[\"versions:1.9\",\"versions:1.8.2\",\"versions:1.15\"],[\"categories:neoforge\"],[\"project_type:mod\"]]", false,
		}, {
			expectedData{"forge",
				"modpack",
				[]string{"1.16.3", "1.13.2", "1.8"}}, "[[\"versions:1.16.3\",\"versions:1.13.2\",\"versions:1.8\"],[\"categories:forge\"],[\"project_type:modpack\"]]", false,
		},
	}

	for _, tt := range test {
		result := FacetsBuilder(tt.given.versions, tt.given.loader, tt.given.projectType)
		if result != tt.expected {
			t.Errorf("got %s expected %s", result, tt.expected)
		}
	}
}

func TestISO8691Format(t *testing.T) {
	tests := []testData[string, string]{
		{"2024-12-19T17:51:48.102945Z", "Dec 19, 2024, 05:51 PM", false},
		{"2023-01-01T00:00:00Z", "Jan 1, 2023, 12:00 AM", false},
		{"1999-07-04T12:30:15Z", "Jul 4, 1999, 12:30 PM", false},
		{"2020-02-29T23:59:59Z", "Feb 29, 2020, 11:59 PM", false},
		{"obviously not a date", "", true},
	}

	for _, tt := range tests {
		result, err := Iso8601Format(tt.given)

		if (err != nil) != tt.err {
			t.Error("got an error, didn't expected one")
		}

		if result != tt.expected {
			t.Errorf("got %s, expected %s", result, tt.expected)
		}
	}
}
