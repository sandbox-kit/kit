package sandbox

import (
	fmt "fmt"
	math "math"
)

// metadataBudget tracks remaining node and payload byte allowance while
// building declarative metadata.
type metadataBudget struct {
	nodes int
	bytes int
}

func metadataError(message string) error {
	return NewError(&ErrorInfo{
		Kind:    ErrorKindInvalidArgument,
		Field:   InfoFieldProviderMetadata,
		Message: message,
	}, nil)
}

// copyMetadataValue checks the depth, node, and payload budgets and builds the
// protobuf value in the same walk. Byte slices are copied once.
func copyMetadataValue(value any, depth int, budget *metadataBudget) (*MetadataValue, error) {
	if depth > 64 {
		return nil, metadataError("sandbox-kit: metadata exceeds maximum depth or contains a cycle")
	}
	if budget.nodes <= 0 {
		return nil, metadataError("sandbox-kit: metadata exceeds node budget")
	}
	budget.nodes--
	switch v := value.(type) {
	case nil:
		return &MetadataValue{}, nil
	case bool:
		return &MetadataValue{Value: &MetadataBoolean{Boolean: v}}, nil
	case string:
		if len(v) > budget.bytes {
			return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
		}
		budget.bytes -= len(v)
		return &MetadataValue{Value: &MetadataText{Text: v}}, nil
	case int:
		return signedMetadata(int64(v)), nil
	case int8:
		return signedMetadata(int64(v)), nil
	case int16:
		return signedMetadata(int64(v)), nil
	case int32:
		return signedMetadata(int64(v)), nil
	case int64:
		return signedMetadata(v), nil
	case uint:
		return unsignedMetadata(uint64(v)), nil
	case uint8:
		return unsignedMetadata(uint64(v)), nil
	case uint16:
		return unsignedMetadata(uint64(v)), nil
	case uint32:
		return unsignedMetadata(uint64(v)), nil
	case uint64:
		return unsignedMetadata(v), nil
	case float32:
		number := float64(v)
		if math.IsNaN(number) || math.IsInf(number, 0) {
			return nil, metadataError("sandbox-kit: metadata requires finite numbers")
		}
		return &MetadataValue{Value: &MetadataNumber{Number: number}}, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, metadataError("sandbox-kit: metadata requires finite numbers")
		}
		return &MetadataValue{Value: &MetadataNumber{Number: v}}, nil
	case []byte:
		if v == nil {
			return &MetadataValue{Value: &MetadataBinary{}}, nil
		}
		if len(v) > budget.bytes {
			return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
		}
		budget.bytes -= len(v)
		return &MetadataValue{Value: &MetadataBinary{Binary: append([]byte{}, v...)}}, nil
	case []string:
		if v == nil {
			return &MetadataValue{Value: &MetadataListValue{List: &MetadataList{Values: []*MetadataValue{}}}}, nil
		}
		if len(v) > budget.nodes {
			return nil, metadataError("sandbox-kit: metadata exceeds node budget")
		}
		budget.nodes -= len(v)
		list := &MetadataList{Values: make([]*MetadataValue, len(v))}
		for i, item := range v {
			if len(item) > budget.bytes {
				return nil, metadataError("sandbox-kit: metadata exceeds byte budget")
			}
			budget.bytes -= len(item)
			list.Values[i] = &MetadataValue{Value: &MetadataText{Text: item}}
		}
		return &MetadataValue{Value: &MetadataListValue{List: list}}, nil
	case []any:
		if v == nil {
			return &MetadataValue{Value: &MetadataListValue{List: &MetadataList{Values: []*MetadataValue{}}}}, nil
		}
		if len(v) > budget.nodes {
			return nil, metadataError("sandbox-kit: metadata exceeds node budget")
		}
		list := &MetadataList{Values: make([]*MetadataValue, len(v))}
		for i, item := range v {
			copied, err := copyMetadataValue(item, depth+1, budget)
			if err != nil {
				return nil, err
			}
			list.Values[i] = copied
		}
		return &MetadataValue{Value: &MetadataListValue{List: list}}, nil
	case map[string]any:
		if v == nil {
			return &MetadataValue{Value: &MetadataObjectValue{}}, nil
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
		object := &MetadataObject{Values: make(map[string]*MetadataValue, len(v))}
		for key, item := range v {
			copied, err := copyMetadataValue(item, depth+1, budget)
			if err != nil {
				return nil, err
			}
			object.Values[key] = copied
		}
		return &MetadataValue{Value: &MetadataObjectValue{Object: object}}, nil
	default:
		return nil, metadataError(fmt.Sprintf("sandbox-kit: unsupported metadata value type %T", value))
	}
}

func signedMetadata(value int64) *MetadataValue {
	return &MetadataValue{Value: &MetadataSignedInteger{SignedInteger: value}}
}

func unsignedMetadata(value uint64) *MetadataValue {
	return &MetadataValue{Value: &MetadataUnsignedInteger{UnsignedInteger: value}}
}
