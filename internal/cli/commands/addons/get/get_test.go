package get

import "testing"

func TestValidate(t *testing.T) {
	test := []struct {
		loader string
		err    bool
	}{
		{"fabric", false},
		{"forge", false},
		{"", true},
		{"invalid", true},
	}

	cmd := NewCommand()

	for _, tt := range test {
		t.Run(tt.loader, func(t *testing.T) {
			if err := cmd.Flags().Set("loader", tt.loader); err != nil {
				t.Fatal(err)
			}
			err := validate(cmd, []string{"test"})
			if (err != nil) != tt.err {
				t.Errorf("got error %v, expected error %v", err, tt.err)
			}
		})
	}
}
