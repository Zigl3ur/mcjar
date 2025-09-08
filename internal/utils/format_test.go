package utils

import (
	"slices"
	"testing"
)

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
	tests := []struct {
		given    []string
		expected map[string][]string
	}{
		{[]string{"1.21.6", "1.8.9", "23w6a", "1.12.2", "1.7.2_pre4", "23w8b"}, map[string][]string{"versions": {"1.21.6", "1.12.2", "1.8.9"}, "snapshots": {"23w6a", "1.7.2_pre4", "23w8b"}}},
		{[]string{"1.9", "1.10.2", "1.4.3"}, map[string][]string{"versions": {"1.10.2", "1.9", "1.4.3"}, "snapshots": {}}},
		{[]string{"1.21.6", "23w6a", "1.21.6"}, map[string][]string{"versions": {"1.21.6", "1.21.6"}, "snapshots": {"23w6a"}}},
		{[]string{"21E", " "}, map[string][]string{"versions": {}, "snapshots": {"21E", " "}}},
		{[]string{""}, map[string][]string{"versions": {""}, "snapshots": {}}},
	}

	for _, tt := range tests {
		result := SortMcVersions(tt.given)
		if !slices.Equal(result["versions"], tt.expected["versions"]) || !slices.Equal(result["snapshots"], tt.expected["snapshots"]) {
			t.Errorf("got %s, expected %s", result, tt.expected)
		}
	}

}

func TestHumanizeByte(t *testing.T) {
	tests := []struct {
		given    int64
		expected string
	}{
		{1, "1 B"},
		{1024, "1.0 KiB"},
		{1048576, "1.0 MiB"},
		{1073741824, "1.0 GiB"},
		{1099511627776, "1.0 TiB"},
		{1125899906842624, "1.0 PiB"},
		{512, "512 B"},
		{000000, "0 B"},
		{2097152, "2.0 MiB"},
		{1610612736, "1.5 GiB"},
		{999, "999 B"},
	}

	for _, tt := range tests {
		result := humanizeByte(tt.given)
		if result != tt.expected {
			t.Errorf("got %s, expected %s", result, tt.expected)
		}
	}
}

func TestISO8601Format(t *testing.T) {
	tests := []struct {
		given    string
		expected string
	}{
		{"2024-12-19T17:51:48.102945Z", "Dec 19, 2024, 05:51 PM"},
		{"2023-01-01T00:00:00Z", "Jan 1, 2023, 12:00 AM"},
		{"1999-07-04T12:30:15Z", "Jul 4, 1999, 12:30 PM"},
		{"2020-02-29T23:59:59Z", "Feb 29, 2020, 11:59 PM"},
		{"obviously not a date", "obviously not a date"},
	}

	for _, tt := range tests {
		result := Iso8601Format(tt.given)

		if result != tt.expected {
			t.Errorf("got %s, expected %s", result, tt.expected)
		}
	}
}
