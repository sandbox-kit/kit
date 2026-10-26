package daytona

import (
	"reflect"

	"google.golang.org/protobuf/proto"
)

// cloneMessage copies a protobuf message without copying its embedded mutex.
// proto.Clone drops empty slices and maps; those are restored so an explicit
// empty collection still counts as set.
func cloneMessage[T proto.Message](src T) T {
	value := reflect.ValueOf(src)
	if !value.IsValid() || (value.Kind() == reflect.Pointer && value.IsNil()) {
		var zero T
		return zero
	}
	dst := proto.Clone(src).(T)
	preserveEmpty(value, reflect.ValueOf(dst))
	return dst
}

func preserveEmpty(src, dst reflect.Value) {
	if src.Kind() == reflect.Pointer {
		if src.IsNil() || dst.Kind() != reflect.Pointer || dst.IsNil() {
			return
		}
		src, dst = src.Elem(), dst.Elem()
	}
	if src.Kind() != reflect.Struct || dst.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < src.NumField(); i++ {
		if src.Type().Field(i).PkgPath != "" {
			continue
		}
		sv, dv := src.Field(i), dst.Field(i)
		switch sv.Kind() {
		case reflect.Slice:
			if !sv.IsNil() && dv.IsNil() {
				dv.Set(reflect.MakeSlice(sv.Type(), 0, 0))
			}
		case reflect.Map:
			if !sv.IsNil() && dv.IsNil() {
				dv.Set(reflect.MakeMap(sv.Type()))
			}
		case reflect.Pointer:
			if sv.IsNil() || dv.IsNil() || sv.Elem().Kind() != reflect.Struct {
				continue
			}
			preserveEmpty(sv, dv)
		case reflect.Struct:
			preserveEmpty(sv, dv)
		}
	}
}
