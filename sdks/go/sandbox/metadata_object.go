package sandbox

import (
	"fmt"
	"math"
)

// MetadataObjectFromMap builds the metadata object accepted on a create request.
// Keys are strings. Values are booleans, strings, integers, finite floats, byte slices, lists, and nested maps.
// A nil map returns a nil object. A depth above 64, more than 16384 values, a non-finite number, or more than 1 MiB of text and bytes returns *Error with ErrorKindInvalidArgument.
func MetadataObjectFromMap(source map[string]any) (*MetadataObject, error) {
	if source == nil {
		return nil, nil
	}
	copied, err := cloneMetadata(source)
	if err != nil {
		return nil, err
	}
	return metadataObjectFromCopiedMap(copied)
}

func metadataObjectFromCopiedMap(source map[string]any) (*MetadataObject, error) {
	if source == nil {
		return nil, nil
	}
	out := &MetadataObject{Values: make(map[string]*MetadataValue, len(source))}
	for key, value := range source {
		converted, err := metadataValueFromAny(value)
		if err != nil {
			return nil, err
		}
		out.Values[key] = converted
	}
	return out, nil
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

func metadataValueFromAny(value any) (*MetadataValue, error) {
	switch v := value.(type) {
	case nil:
		return &MetadataValue{}, nil
	case bool:
		return &MetadataValue{Value: &MetadataBoolean{Boolean: v}}, nil
	case string:
		return &MetadataValue{Value: &MetadataText{Text: v}}, nil
	case int:
		return &MetadataValue{Value: &MetadataSignedInteger{SignedInteger: int64(v)}}, nil
	case int8:
		return &MetadataValue{Value: &MetadataSignedInteger{SignedInteger: int64(v)}}, nil
	case int16:
		return &MetadataValue{Value: &MetadataSignedInteger{SignedInteger: int64(v)}}, nil
	case int32:
		return &MetadataValue{Value: &MetadataSignedInteger{SignedInteger: int64(v)}}, nil
	case int64:
		return &MetadataValue{Value: &MetadataSignedInteger{SignedInteger: v}}, nil
	case uint:
		return &MetadataValue{Value: &MetadataUnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint8:
		return &MetadataValue{Value: &MetadataUnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint16:
		return &MetadataValue{Value: &MetadataUnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint32:
		return &MetadataValue{Value: &MetadataUnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint64:
		return &MetadataValue{Value: &MetadataUnsignedInteger{UnsignedInteger: v}}, nil
	case float32:
		return &MetadataValue{Value: &MetadataNumber{Number: float64(v)}}, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("sandbox-kit: metadata numbers must be finite")
		}
		return &MetadataValue{Value: &MetadataNumber{Number: v}}, nil
	case []byte:
		return &MetadataValue{Value: &MetadataBinary{Binary: append([]byte(nil), v...)}}, nil
	case []string:
		list := &MetadataList{Values: make([]*MetadataValue, 0, len(v))}
		for _, item := range v {
			list.Values = append(list.Values, &MetadataValue{Value: &MetadataText{Text: item}})
		}
		return &MetadataValue{Value: &MetadataListValue{List: list}}, nil
	case []any:
		list := &MetadataList{Values: make([]*MetadataValue, 0, len(v))}
		for _, item := range v {
			converted, err := metadataValueFromAny(item)
			if err != nil {
				return nil, err
			}
			list.Values = append(list.Values, converted)
		}
		return &MetadataValue{Value: &MetadataListValue{List: list}}, nil
	case map[string]any:
		object, err := metadataObjectFromCopiedMap(v)
		if err != nil {
			return nil, err
		}
		return &MetadataValue{Value: &MetadataObjectValue{Object: object}}, nil
	default:
		return nil, fmt.Errorf("sandbox-kit: unsupported metadata value %T", value)
	}
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
