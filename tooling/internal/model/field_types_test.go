package model

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

func TestFieldProjectionPreservesWidthAndGetterPresence(t *testing.T) {
	schema := NewSchema([]protoreflect.FileDescriptor{fixture(t)})
	path, err := schema.Resolve("Resources", "cpu_cores")
	if err != nil {
		t.Fatal(err)
	}
	stored, err := FieldType(path.Fields[0], false)
	if err != nil {
		t.Fatal(err)
	}
	getter, err := FieldType(path.Fields[0], true)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Builtin != "number" || !stored.Optional || getter.Builtin != "number" || getter.Optional {
		t.Fatal("scalar shape or presence lost")
	}
	path, err = schema.Resolve("CreateOptions", "resources")
	if err != nil {
		t.Fatal(err)
	}
	object, err := FieldType(path.Fields[0], true)
	if err != nil {
		t.Fatal(err)
	}
	if object.Ref != "schema.kit.sandbox.v1.Resources" || !object.Reference {
		t.Fatal("message identity lost")
	}
}
