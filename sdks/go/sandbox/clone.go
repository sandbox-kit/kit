package sandbox

import (
	"encoding/json"
	"fmt"
)

// cloneData copies declarative configuration and response metadata, never SDK objects.
// It preserves optional values and ownership boundaries.
func cloneData[T any](value *T) (*T, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("sandbox-kit: configuration/metadata must be declarative: %w", err)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
