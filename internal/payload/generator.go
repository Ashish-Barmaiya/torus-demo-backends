package payload

import (
	"encoding/json"
	"fmt"
)

const defaultFiller = "torus-demo-payload"

type Envelope struct {
	Data    any    `json:"data"`
	Meta    any    `json:"meta"`
	Payload string `json:"payload,omitempty"`
}

func GenerateJSON(target Size) ([]byte, error) {
	return GenerateJSONWithData(
		target,
		nil,
		nil,
	)
}

func GenerateJSONWithData(
	target Size,
	data any,
	meta any,
) ([]byte, error) {
	if !target.Valid() {
		return nil, fmt.Errorf("unsupported payload size: %d", target)
	}

	if target == SizeEmpty {
		return []byte{}, nil
	}

	envelope := Envelope{
		Data: data,
		Meta: meta,
	}

	// Build the response structure using a binary-search-style
	// adjustment because JSON escaping can affect encoded length.
	low := 0
	high := int(target)

	for low <= high {
		mid := (low + high) / 2

		envelope.Payload = makeFiller(mid)

		body, err := json.Marshal(envelope)
		if err != nil {
			return nil, fmt.Errorf("marshal response: %w", err)
		}

		switch {
		case len(body) < int(target):
			low = mid + 1

		case len(body) > int(target):
			high = mid - 1

		default:
			return body, nil
		}
	}

	// Need the exact size, but a JSON string can require escaped bytes.
	// Find the closest valid result and refine it.
	best := []byte(nil)
	bestDistance := int(target)

	for fillerLength := max(0, high-4); fillerLength <= min(int(target), high+4); fillerLength++ {
		envelope.Payload = makeFiller(fillerLength)

		body, err := json.Marshal(envelope)
		if err != nil {
			return nil, fmt.Errorf("marshal response: %w", err)
		}

		distance := abs(len(body) - int(target))

		if distance < bestDistance {
			best = body
			bestDistance = distance
		}

		if distance == 0 {
			return body, nil
		}
	}

	if best == nil {
		return nil, fmt.Errorf(
			"unable to generate payload near target size %d",
			target,
		)
	}

	return best, nil
}

func makeFiller(length int) string {
	if length <= 0 {
		return ""
	}

	pattern := []byte(defaultFiller)
	filler := make([]byte, length)

	for i := range filler {
		filler[i] = pattern[i%len(pattern)]
	}

	return string(filler)
}

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
