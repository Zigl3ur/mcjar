package list

import "testing"

func TestValidate(t *testing.T) {

	test := []struct {
		serverType string
		err        bool
	}{
		{"paper", false},
		{"invalid", true},
	}

	cmd := NewCommand()

	for _, tt := range test {
		t.Run(tt.serverType, func(t *testing.T) {
			if err := cmd.Flags().Set("type", tt.serverType); err != nil {
				t.Fatal(err)
			}
			err := validate(cmd, nil)
			if (err != nil) != tt.err {
				t.Errorf("got error %v, expected error %v", err, tt.err)
			}
		})
	}
}
