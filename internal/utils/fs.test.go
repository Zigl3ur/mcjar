package utils

import (
	"os"
	"testing"
)

func TestGetPath(t *testing.T) {
	test := []struct {
		given string
		err   bool
	}{
		{"go", false},
		{"not-an-existing-command", true},
	}

	for _, tt := range test {
		_, err := GetPath(tt.given)
		if (err != nil) != tt.err {
			t.Errorf("got error %v, expected error: %v", err, tt.err)
		}
	}
}

func TestCheckDir(t *testing.T) {
	test := []struct {
		given string
		err   bool
	}{
		{"./testdata/checkdir", false},
		{"/root/forbidden-dir", true},
	}

	for _, tt := range test {
		err := CheckDir(tt.given)
		if (err != nil) != tt.err {
			t.Errorf("got error %v, expected error: %v", err, tt.err)
		}
	}

	err := os.RemoveAll("testdata/")
	if err != nil {
		t.Errorf("failed to clean up test directory: %v", err)
	}
}
