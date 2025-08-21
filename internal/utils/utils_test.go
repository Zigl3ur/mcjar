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
