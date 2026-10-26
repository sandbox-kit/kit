package sandbox

import (
	fmt "fmt"
	math "math"
)

// metadataBudget tracks remaining node and payload byte allowance while
// cloning declarative metadata.
type metadataBudget struct {
	nodes int
	bytes int
}

// cloneMetadata returns an owned copy of declarative metadata. A nil map stays
// nil. SDK objects are rejected.
func cloneMetadata(source map[string]any) (map[string]any, error) {
	if source == nil {
		return nil, nil
	}
	budget := &metadataBudget{nodes: 16384, bytes: 1048576}
	value, err := copyMetadataValue(source, 0, budget)
	if err != nil {
		return nil, err
	}
	return value.(map[string]any), nil
}

func metadataError(message string) error {
	return NewError(&ErrorInfo{
		Kind:    ErrorKindInvalidArgument,
		Field:   InfoFieldProviderMetadata,
		Message: message,
	}, nil)
}

// copyMetadataValue copies metadata within the depth, node, and payload byte
// budgets.
func copyMetadataValue(value any, depth int, budget *metadataBudget) (any, error) {
	if depth > 64 {
		return nil, metadataError("sandbox-kit: metadata exceeds maximum depth or contains a cycle")
	}
	if budget.nodes <= 0 {
		return nil, metadataError("sandbox-kit: metadata exceeds node budget")
	}
	budget.nodes--
	switch v := value.(type) {
	case nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v, nil
	case string:
		if len(v) > budget.bytes {
			return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
		}
		budget.bytes -= len(v)
		return v, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, metadataError("sandbox-kit: metadata requires finite numbers")
		}
		return v, nil
	case float32:
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			return nil, metadataError("sandbox-kit: metadata requires finite numbers")
		}
		return v, nil
	case []byte:
		if v == nil {
			return []byte(nil), nil
		}
		if len(v) > budget.bytes {
			return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
		}
		budget.bytes -= len(v)
		return append([]byte{}, v...), nil
	case []string:
		if v == nil {
			return []string(nil), nil
		}
		if len(v) > budget.nodes {
			return nil, metadataError("sandbox-kit: metadata exceeds node budget")
		}
		budget.nodes -= len(v)
		for _, item := range v {
			if len(item) > budget.bytes {
				return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
			}
			budget.bytes -= len(item)
		}
		return append([]string{}, v...), nil
	case []any:
		if v == nil {
			return []any(nil), nil
		}
		if len(v) > budget.nodes {
			return nil, metadataError("sandbox-kit: metadata exceeds node budget")
		}
		out := make([]any, len(v))
		for i, item := range v {
			copied, err := copyMetadataValue(item, depth+1, budget)
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
		if len(v) > budget.nodes {
			return nil, metadataError("sandbox-kit: metadata exceeds node budget")
		}
		for key := range v {
			if len(key) > budget.bytes {
				return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
			}
			budget.bytes -= len(key)
		}
		out := make(map[string]any, len(v))
		for key, item := range v {
			copied, err := copyMetadataValue(item, depth+1, budget)
			if err != nil {
				return nil, err
			}
			out[key] = copied
		}
		return out, nil
	default:
		return nil, metadataError(fmt.Sprintf("sandbox-kit: unsupported metadata value type %T", value))
	}
}
