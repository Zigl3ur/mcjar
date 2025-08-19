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
	tests := []testData[string, [3]int]{
		{"1.8.9", [3]int{1, 8, 9}, false},
		{"1.12.2", [3]int{1, 12, 2}, false},
		{"1.10", [3]int{1, 10, 0}, false},
		{"24w5a", [3]int{0, 0, 0}, true},
		{"24w5a.12.ZAE", [3]int{0, 12, 0}, false},
		{"", [3]int{0, 0, 0}, true},
	}

	for _, tt := range tests {
		result, err := mcVersionParser(tt.given)
		if (err != nil) != tt.err {
			t.Errorf("McVersionParser(%s): got error = %t, want error = %t", tt.given, err != nil, tt.err)
		}

		if result != tt.expected {
			t.Errorf("got %d, expected %d", result, tt.expected)
		}
	}

}

func TestSortMcVersions(t *testing.T) {
	tests := []testData[[]string, []string]{
		{[]string{"1.21.6", "1.8.9", "23w6a", "1.12.2", "1.7.2", "23w8b"}, []string{
			"1.7.2", "1.8.9", "23w6a", "1.12.2", "1.21.6", "23w8b"}, false},
		{[]string{"1.9", "1.10.2", "1.4.3"}, []string{"1.4.3", "1.9", "1.10.2"}, false},
		{[]string{"1.21.6", "23w6a", "1.21.6"}, []string{"1.21.6", "23w6a", "1.21.6"}, false},
		{[]string{"21E", " "}, []string{"21E", " "}, false},
		{[]string{""}, []string{""}, false},
	}

	for _, tt := range tests {
		result := SortMcVersions(tt.given)
		if !slices.Equal(result, tt.expected) {
			t.Errorf("got %s expected %s", result, tt.expected)
		}
	}

}

func TestHumanizeByte(t *testing.T) {
	tests := []testData[int64, string]{
		{1, "1 B", false},
		{1000, "1.0 kB", false},
		{1000000, "1.0 MB", false},
		{1000000000, "1.0 GB", false},
		{1000000000000, "1.0 TB", false},
		{1000000000000000, "1.0 PB", false},
		{512, "512 B", false},
		{000000, "0 B", false},
		{2097152, "2.1 MB", false},
		{1610612736, "1.6 GB", false},
		{999, "999 B", false},
	}

	for _, tt := range tests {
		result := humanizeByte(tt.given)
		if result != tt.expected {
			t.Errorf("got %s expected %s", result, tt.expected)
		}
	}
}
