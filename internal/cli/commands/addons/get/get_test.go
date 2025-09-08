package get

import "testing"

func TestValidate(t *testing.T) {
	test := []struct {
		loader string
		err    bool
	}{
		{"fabric", false},
		{"forge", false},
		{"", false},
		{"invalid", true},
	}

	cmd := NewCommand()

	for _, tt := range test {
		t.Run(tt.loader, func(t *testing.T) {
			cmd.Flags().Set("loader", tt.loader)
			err := validate(cmd, []string{"test"})
			if (err != nil) != tt.err {
				t.Errorf("got error %v, expected error %v", err, tt.err)
			}
		})
	}
}
