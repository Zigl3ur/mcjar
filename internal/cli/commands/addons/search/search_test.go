package search

import "testing"

func TestValidate(t *testing.T) {
	test := []struct {
		loader     string
		addonsType string
		err        bool
	}{
		{"fabric", "mod", false},
		{"forge", "mod", false},
		{"invalid", "mod", true},
		{"fabric", "invalid", true},
	}

	cmd := NewCommand()

	for _, tt := range test {
		t.Run(tt.loader+tt.addonsType, func(t *testing.T) {
			cmd.Flags().Set("loader", tt.loader)
			cmd.Flags().Set("type", tt.addonsType)
			err := validate(cmd, []string{"test"})
			if (err != nil) != tt.err {
				t.Errorf("got error %v, expected error %v", err, tt.err)
			}
		})
	}

}
