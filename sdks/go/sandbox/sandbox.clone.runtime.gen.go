// Code generated from contracts.yaml and errors.yaml by the Go emitter. DO NOT EDIT.
package sandbox

import (
	fmt "fmt"
	math "math"
)

// cloneMetadata copies declarative metadata, never SDK objects.
// It preserves optional values and ownership boundaries.
func cloneMetadata(source map[string]any) (map[string]any, error) {
	if source == nil {
		return nil, nil
	}
	value, err := copyMetadataValue(source, 0)
	if err != nil {
		return nil, err
	}
	return value.(map[string]any), nil
}

// Metadata is declarative: finite primitive values, bytes, lists, and string-key
// objects. Depth bounds reject cycles without retaining caller-owned objects.
func copyMetadataValue(value any, depth int) (any, error) {
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: metadata exceeds maximum depth or contains a cycle")
	}
	switch v := value.(type) {
	case nil, bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("sandbox-kit: metadata requires finite numbers")
		}
		return v, nil
	case float32:
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			return nil, fmt.Errorf("sandbox-kit: metadata requires finite numbers")
		}
		return v, nil
	case []byte:
		if v == nil {
			return []byte(nil), nil
		}
		return append([]byte{}, v...), nil
	case []string:
		if v == nil {
			return []string(nil), nil
		}
		return append([]string{}, v...), nil
	case []any:
		if v == nil {
			return []any(nil), nil
		}
		out := make([]any, len(v))
		for i, item := range v {
			copied, err := copyMetadataValue(item, depth+1)
			if err != nil {
				return nil, err
			}
			out[i] = copied
		}
		return out, nil
	case map[string]any:
		if v == nil {
			return map[string]any(nil), nil
		}
		out := make(map[string]any, len(v))
		for key, item := range v {
			copied, err := copyMetadataValue(item, depth+1)
			if err != nil {
				return nil, err
			}
			out[key] = copied
		}
		return out, nil
	default:
		return nil, fmt.Errorf("sandbox-kit: unsupported metadata value type %T", value)
	}
}
