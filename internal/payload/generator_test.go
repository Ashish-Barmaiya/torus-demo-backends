package payload

import (
	"encoding/json"
	"testing"
)

func TestGenerateJSON(t *testing.T) {
	tests := []Size{
		Size1KB,
		Size16KB,
		Size64KB,
		Size256KB,
		Size1MB,
		Size4MB,
	}

	for _, size := range tests {
		t.Run(size.String(), func(t *testing.T) {
			body, err := GenerateJSON(size)
			if err != nil {
				t.Fatalf("GenerateJSON() error: %v", err)
			}

			if len(body) != int(size) {
				t.Fatalf(
					"expected %d bytes, got %d",
					size,
					len(body),
				)
			}

			var decoded any

			if err := json.Unmarshal(body, &decoded); err != nil {
				t.Fatalf("generated payload is not valid JSON: %v", err)
			}
		})
	}
}

func TestGenerateJSONEmpty(t *testing.T) {
	body, err := GenerateJSON(SizeEmpty)
	if err != nil {
		t.Fatalf("GenerateJSON() error: %v", err)
	}

	if len(body) != 0 {
		t.Fatalf("expected empty payload, got %d bytes", len(body))
	}
}

func TestGenerateJSONUnsupportedSize(t *testing.T) {
	tests := []Size{
		1,
		2048,
		8 << 20,
		-1,
	}

	for _, size := range tests {
		t.Run(size.String(), func(t *testing.T) {
			body, err := GenerateJSON(size)

			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if body != nil {
				t.Fatalf("expected nil body, got %d bytes", len(body))
			}
		})
	}
}

func TestSizeValid(t *testing.T) {
	valid := []Size{
		SizeEmpty,
		Size1KB,
		Size16KB,
		Size64KB,
		Size256KB,
		Size1MB,
		Size4MB,
	}

	for _, size := range valid {
		if !size.Valid() {
			t.Fatalf("expected %s to be valid", size)
		}
	}

	invalid := []Size{
		1,
		2 << 10,
		8 << 20,
		-1,
	}

	for _, size := range invalid {
		if size.Valid() {
			t.Fatalf("expected %s to be invalid", size)
		}
	}
}
