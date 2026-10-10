// Package model compiles shared descriptors and portable specifications into a
// language-neutral generation model. It has no protogen or SDK dependencies.
package model

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

// Schema indexes full protobuf identities, without target-language names.
type Schema struct {
	messages map[protoreflect.FullName]protoreflect.MessageDescriptor
}

func NewSchema(files []protoreflect.FileDescriptor) *Schema {
	s := &Schema{messages: map[protoreflect.FullName]protoreflect.MessageDescriptor{}}
	var visit func(protoreflect.MessageDescriptors)
	visit = func(messages protoreflect.MessageDescriptors) {
		for i := 0; i < messages.Len(); i++ {
			m := messages.Get(i)
			s.messages[m.FullName()] = m
			visit(m.Messages())
		}
	}
	for _, file := range files {
		visit(file.Messages())
	}
	return s
}

func (s *Schema) Message(name string) (protoreflect.MessageDescriptor, error) {
	if !strings.Contains(name, ".") {
		name = "kit.sandbox.v1." + name
	}
	m := s.messages[protoreflect.FullName(name)]
	if m == nil {
		return nil, fmt.Errorf("unknown shared message %s", name)
	}
	return m, nil
}

// FieldPath preserves canonical names, numeric identities, and presence.
type FieldPath struct {
	Root   protoreflect.FullName
	Path   string
	Fields []protoreflect.FieldDescriptor
}

func (s *Schema) Resolve(root, path string) (FieldPath, error) {
	m, err := s.Message(root)
	if err != nil {
		return FieldPath{}, err
	}
	result := FieldPath{Root: m.FullName(), Path: path}
	parts := strings.Split(path, ".")
	for i, part := range parts {
		field := m.Fields().ByName(protoreflect.Name(part))
		if field == nil {
			return FieldPath{}, fmt.Errorf("unknown path %s.%s", root, path)
		}
		result.Fields = append(result.Fields, field)
		if i < len(parts)-1 {
			if field.Message() == nil || field.IsMap() {
				return FieldPath{}, fmt.Errorf("path %s traverses a scalar or map", path)
			}
			name := string(field.Message().FullName())
			if strings.HasPrefix(name, "google.protobuf.") || name == "kit.sandbox.v1.MetadataObject" {
				return FieldPath{}, fmt.Errorf("path %s traverses a native data leaf", path)
			}
			m = field.Message()
		}
	}
	return result, nil
}
