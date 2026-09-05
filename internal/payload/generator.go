package payload

import (
	"encoding/json"
	"fmt"
)

const (
	defaultFiller = "torus-demo-payload"
)

type Envelope struct {
	Data string `json:"data"`
}

func GenerateJSON(target Size) ([]byte, error) {
	if !target.Valid() {
		return nil, fmt.Errorf("unsupported payload size: %d", target)
	}

	if target == SizeEmpty {
		return []byte{}, nil
	}

	// Generate a JSON document whose encoded size is exactly target bytes.
	//
	// The JSON object has a fixed prefix/suffix and a variable-length string.
	prefix := []byte(`{"data":"`)
	suffix := []byte(`"}`)

	if target <= Size(len(prefix)+len(suffix)) {
		return nil, fmt.Errorf(
			"target size %d is too small for JSON envelope",
			target,
		)
	}

	fillerLength := int(target) - len(prefix) - len(suffix)

	filler := make([]byte, fillerLength)

	pattern := []byte(defaultFiller)

	for i := range filler {
		filler[i] = pattern[i%len(pattern)]
	}

	body := make([]byte, 0, int(target))
	body = append(body, prefix...)
	body = append(body, filler...)
	body = append(body, suffix...)

	if len(body) != int(target) {
		return nil, fmt.Errorf(
			"generated payload size mismatch: got %d, want %d",
			len(body),
			target,
		)
	}

	// Verify that the result is actually valid JSON.
	var decoded Envelope

	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, fmt.Errorf("generated invalid JSON: %w", err)
	}

	return body, nil
}
