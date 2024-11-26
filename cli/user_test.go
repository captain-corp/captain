package cli

import (
	"os"
	"testing"
)

func TestGetValidInput(t *testing.T) {
	// Save original stdin
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	tests := []struct {
		name      string
		input     string
		validator func(string) error
		want      string
	}{
		{
			name:      "Valid input",
			input:     "test\n",
			validator: validateFirstName,
			want:      "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe and make it stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write test input
			w.Write([]byte(tt.input))
			w.Close()

			if got := getValidInput("Test: ", tt.validator); got != tt.want {
				t.Errorf("getValidInput() = %v, want %v", got, tt.want)
			}
		})
	}
}