package model

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sort"
)

// ValidateRules checks semantic field references and constraint vocabulary.
func (s *Schema) ValidateRules(rules spec.Validation) error {
	names := make([]string, 0, len(rules.Messages))
	for name := range rules.Messages {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		message := rules.Messages[name]
		if _, err := s.Message(name); err != nil {
			return err
		}
		fields := make([]string, 0, len(message.Fields))
		for field := range message.Fields {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		for _, field := range fields {
			path, err := s.Resolve(name, field)
			if err != nil {
				return err
			}
			format := message.Fields[field].Format
			if format != "" {
				switch format {
				case "http_url", "secure_endpoint", "image_reference":
				default:
					return fmt.Errorf("%s.%s: unknown format %s", name, field, format)
				}
				if path.Fields[len(path.Fields)-1].Kind() != protoreflect.StringKind {
					return fmt.Errorf("%s.%s: format requires string", name, field)
				}
			}
		}
		for _, constraint := range message.Constraints {
			switch constraint.Op {
			case "at_most_one", "exactly_one", "requires_any", "forbids", "limit", "requires_positive":
			default:
				return fmt.Errorf("%s: unknown constraint %s", name, constraint.Op)
			}
			paths := append([]string{}, constraint.Fields...)
			if constraint.Field != "" {
				paths = append(paths, constraint.Field)
			}
			if constraint.Other != "" {
				paths = append(paths, constraint.Other)
			}
			if constraint.When != nil {
				paths = append(paths, constraint.When.Field)
			}
			for _, path := range paths {
				if _, err := s.Resolve(name, path); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
