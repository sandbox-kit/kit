package sandbox

import (
	"fmt"
)

// MetadataObjectAsMap returns metadata as a Go map of booleans, strings, numbers, byte slices, lists, and nested maps.
// A nil object returns a nil map. Use this to read SandboxInfo metadata after a create or refresh.
func MetadataObjectAsMap(source *MetadataObject) (map[string]any, error) {
	if source == nil {
		return nil, nil
	}
	out := make(map[string]any, len(source.Values))
	for key, value := range source.Values {
		converted, err := metadataValueAsAny(value)
		if err != nil {
			return nil, err
		}
		out[key] = converted
	}
	return out, nil
}

func metadataValueAsAny(value *MetadataValue) (any, error) {
	if value == nil || value.Value == nil {
		return nil, nil
	}
	switch v := value.Value.(type) {
	case *MetadataBoolean:
		return v.Boolean, nil
	case *MetadataText:
		return v.Text, nil
	case *MetadataSignedInteger:
		return v.SignedInteger, nil
	case *MetadataUnsignedInteger:
		return v.UnsignedInteger, nil
	case *MetadataNumber:
		return v.Number, nil
	case *MetadataBinary:
		return append([]byte(nil), v.Binary...), nil
	case *MetadataListValue:
		out := make([]any, 0, len(v.List.GetValues()))
		for _, item := range v.List.GetValues() {
			converted, err := metadataValueAsAny(item)
			if err != nil {
				return nil, err
			}
			out = append(out, converted)
		}
		return out, nil
	case *MetadataObjectValue:
		return MetadataObjectAsMap(v.Object)
	default:
		return nil, fmt.Errorf("sandbox-kit: unsupported metadata value")
	}
}
