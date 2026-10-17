package sandbox

import (
	"fmt"
	"math"
)

// MetadataObjectFromMap converts declarative map metadata into the protobuf object.
// Conversion applies the same depth, node, and byte budgets as cloneMetadata.
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

// MetadataObjectAsMap converts a protobuf metadata object into a declarative map.
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
		return &MetadataValue{Value: &MetadataValue_Boolean{Boolean: v}}, nil
	case string:
		return &MetadataValue{Value: &MetadataValue_Text{Text: v}}, nil
	case int:
		return &MetadataValue{Value: &MetadataValue_SignedInteger{SignedInteger: int64(v)}}, nil
	case int8:
		return &MetadataValue{Value: &MetadataValue_SignedInteger{SignedInteger: int64(v)}}, nil
	case int16:
		return &MetadataValue{Value: &MetadataValue_SignedInteger{SignedInteger: int64(v)}}, nil
	case int32:
		return &MetadataValue{Value: &MetadataValue_SignedInteger{SignedInteger: int64(v)}}, nil
	case int64:
		return &MetadataValue{Value: &MetadataValue_SignedInteger{SignedInteger: v}}, nil
	case uint:
		return &MetadataValue{Value: &MetadataValue_UnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint8:
		return &MetadataValue{Value: &MetadataValue_UnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint16:
		return &MetadataValue{Value: &MetadataValue_UnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint32:
		return &MetadataValue{Value: &MetadataValue_UnsignedInteger{UnsignedInteger: uint64(v)}}, nil
	case uint64:
		return &MetadataValue{Value: &MetadataValue_UnsignedInteger{UnsignedInteger: v}}, nil
	case float32:
		return &MetadataValue{Value: &MetadataValue_Number{Number: float64(v)}}, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("sandbox-kit: metadata numbers must be finite")
		}
		return &MetadataValue{Value: &MetadataValue_Number{Number: v}}, nil
	case []byte:
		return &MetadataValue{Value: &MetadataValue_Binary{Binary: append([]byte(nil), v...)}}, nil
	case []string:
		list := &MetadataList{Values: make([]*MetadataValue, 0, len(v))}
		for _, item := range v {
			list.Values = append(list.Values, &MetadataValue{Value: &MetadataValue_Text{Text: item}})
		}
		return &MetadataValue{Value: &MetadataValue_List{List: list}}, nil
	case []any:
		list := &MetadataList{Values: make([]*MetadataValue, 0, len(v))}
		for _, item := range v {
			converted, err := metadataValueFromAny(item)
			if err != nil {
				return nil, err
			}
			list.Values = append(list.Values, converted)
		}
		return &MetadataValue{Value: &MetadataValue_List{List: list}}, nil
	case map[string]any:
		object, err := metadataObjectFromCopiedMap(v)
		if err != nil {
			return nil, err
		}
		return &MetadataValue{Value: &MetadataValue_Object{Object: object}}, nil
	default:
		return nil, fmt.Errorf("sandbox-kit: unsupported metadata value %T", value)
	}
}

func metadataValueAsAny(value *MetadataValue) (any, error) {
	if value == nil || value.Value == nil {
		return nil, nil
	}
	switch v := value.Value.(type) {
	case *MetadataValue_Boolean:
		return v.Boolean, nil
	case *MetadataValue_Text:
		return v.Text, nil
	case *MetadataValue_SignedInteger:
		return v.SignedInteger, nil
	case *MetadataValue_UnsignedInteger:
		return v.UnsignedInteger, nil
	case *MetadataValue_Number:
		return v.Number, nil
	case *MetadataValue_Binary:
		return append([]byte(nil), v.Binary...), nil
	case *MetadataValue_List:
		out := make([]any, 0, len(v.List.GetValues()))
		for _, item := range v.List.GetValues() {
			converted, err := metadataValueAsAny(item)
			if err != nil {
				return nil, err
			}
			out = append(out, converted)
		}
		return out, nil
	case *MetadataValue_Object:
		return MetadataObjectAsMap(v.Object)
	default:
		return nil, fmt.Errorf("sandbox-kit: unsupported metadata value")
	}
}
