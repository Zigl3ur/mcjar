package search

import "testing"

func TestValidate(t *testing.T) {
	test := []struct {
		loader     string
		addonsType string
		indexType  string
		err        bool
	}{
		{"fabric", "mod", "relevance", false},
		{"forge", "mod", "newest", false},
		{"invalid", "mod", "noindex", true},
		{"fabric", "invalid", "stillnoindex", true},
	}

	cmd := NewCommand()

	for _, tt := range test {
		t.Run(tt.loader+tt.addonsType, func(t *testing.T) {
			if err := cmd.Flags().Set("loader", tt.loader); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Flags().Set("type", tt.addonsType); err != nil {
				t.Fatal(err)
			}
			if err := cmd.Flags().Set("index", tt.indexType); err != nil {
				t.Fatal(err)
			}
			err := validate(cmd, []string{"test"})
			if (err != nil) != tt.err {
				t.Errorf("got error %v, expected error %v", err, tt.err)
			}
		})
	}

}
