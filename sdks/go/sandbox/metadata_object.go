package sandbox

import (
	"fmt"
)

// MetadataObjectFromMap builds the metadata object accepted on a create request.
// Keys are strings. Values are booleans, strings, integers, finite floats, byte slices, lists, and nested maps.
// A nil map returns a nil object. A depth above 64, more than 16384 values, a non-finite number, or more than 1 MiB of text and bytes returns *Error with ErrorKindInvalidArgument.
func MetadataObjectFromMap(source map[string]any) (*MetadataObject, error) {
	if source == nil {
		return nil, nil
	}
	budget := &metadataBudget{nodes: 16384, bytes: 1048576}
	value, err := copyMetadataValue(source, 0, budget)
	if err != nil {
		return nil, err
	}
	return value.GetObject(), nil
}

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
