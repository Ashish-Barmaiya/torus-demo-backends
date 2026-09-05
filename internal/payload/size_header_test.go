package payload

import "testing"

func TestParseSize(t *testing.T) {
	tests := []struct {
		input string
		want  Size
	}{
		{"", SizeEmpty},
		{"1kb", Size1KB},
		{"16kb", Size16KB},
		{"64kb", Size64KB},
		{"256kb", Size256KB},
		{"1mb", Size1MB},
		{"4mb", Size4MB},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseSize(tt.input)
			if err != nil {
				t.Fatalf("ParseSize() error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}

func TestParseSizeUnsupported(t *testing.T) {
	tests := []string{
		"2kb",
		"8mb",
		"64MB",
		"invalid",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := ParseSize(input)

			if err == nil {
				t.Fatalf("expected error for %q", input)
			}
		})
	}
}
