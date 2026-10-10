package model

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sort"
	"strings"
)

// Kind preserves the contract enum identity; emitters derive native identifiers.
type Kind struct {
	Name   protoreflect.FullName
	Number protoreflect.EnumNumber
}

func (s *Schema) ErrorKind(name string) (Kind, error) {
	m, err := s.Message("ErrorInfo")
	if err != nil {
		return Kind{}, err
	}
	field := m.Fields().ByName("kind")
	if field == nil || field.Enum() == nil {
		return Kind{}, fmt.Errorf("ErrorInfo.kind must be an enum")
	}
	value := field.Enum().Values().ByName(protoreflect.Name("ERROR_KIND_" + strings.ToUpper(name)))
	if value == nil {
		return Kind{}, fmt.Errorf("unknown error kind %s", name)
	}
	return Kind{Name: value.FullName(), Number: value.Number()}, nil
}

// ErrorCase is portable classification. Native rule IDs are bound separately.
type ErrorCase struct {
	Code int
	Rule string
	Kind Kind
}
type ErrorPlan struct {
	Native []ErrorCase
	HTTP   []ErrorCase
	GRPC   []ErrorCase
}
type Mapping struct {
	Rule      spec.Mapping
	Shared    FieldPath
	ErrorPath string
}
type Group struct {
	Rule   spec.Group
	Fields []Mapping
}
type Provider struct {
	Rules  spec.Provider
	Groups []Group
	Errors ErrorPlan
}

func CompileProvider(schema *Schema, rules spec.Provider) (Provider, error) {
	result := Provider{Rules: rules}
	if rules.Client != nil {
		seen := map[string]bool{}
		for _, rule := range rules.Client.Environment {
			path, err := schema.Resolve("Config", rule.Field)
			if err != nil {
				return result, err
			}
			if seen[rule.Field] || len(path.Fields) != 1 || path.Fields[0].Kind() != protoreflect.StringKind || !path.Fields[0].HasPresence() || rule.Default == "" || len(rule.Variables) == 0 {
				return result, fmt.Errorf("invalid environment resolution for %s", rule.Field)
			}
			seen[rule.Field] = true
			for _, variable := range rule.Variables {
				if variable == "" || strings.ContainsAny(variable, "=\x00") {
					return result, fmt.Errorf("invalid environment variable")
				}
			}
		}
	}
	ruleNames := make([]string, 0, len(rules.ErrorRules))
	for rule := range rules.ErrorRules {
		ruleNames = append(ruleNames, rule)
	}
	sort.Strings(ruleNames)
	for _, rule := range ruleNames {
		kind, err := schema.ErrorKind(rules.ErrorRules[rule])
		if err != nil {
			return result, err
		}
		result.Errors.Native = append(result.Errors.Native, ErrorCase{Rule: rule, Kind: kind})
	}
	for _, input := range []struct {
		rules    map[int]string
		target   *[]ErrorCase
		min, max int
	}{{rules.ErrorStatuses, &result.Errors.HTTP, 100, 599}, {rules.ErrorGRPC, &result.Errors.GRPC, 1, 16}} {
		codes := make([]int, 0, len(input.rules))
		for code := range input.rules {
			codes = append(codes, code)
		}
		sort.Ints(codes)
		for _, code := range codes {
			if code < input.min || code > input.max {
				return result, fmt.Errorf("invalid protocol error code %d", code)
			}
			kind, err := schema.ErrorKind(input.rules[code])
			if err != nil {
				return result, err
			}
			*input.target = append(*input.target, ErrorCase{Code: code, Kind: kind})
		}
	}
	for _, binding := range rules.Errors {
		for _, native := range binding.Types {
			if _, ok := rules.ErrorRules[native.Rule]; !ok {
				return result, fmt.Errorf("missing portable error rule %s", native.Rule)
			}
		}
	}
	seen := map[string]bool{}
	for _, group := range rules.Groups {
		if group.Name == "" || seen[group.Name] {
			return result, fmt.Errorf("invalid or duplicate mapping group %s", group.Name)
		}
		seen[group.Name] = true
		shared := group.Source
		if group.Direction == "response" {
			shared = group.Target
		} else if group.Direction != "request" {
			return result, fmt.Errorf("unknown mapping direction %s", group.Direction)
		}
		if len(shared.Bindings) > 0 {
			return result, fmt.Errorf("%s: shared mapping side cannot have native bindings", group.Name)
		}
		if _, err := schema.Message(shared.Name); err != nil {
			return result, err
		}
		compiled := Group{Rule: group}
		if group.Origin != "" && (group.Direction != "response" || group.Target.Name != "SandboxInfo" || (group.Origin != "provider" && group.Origin != "request")) {
			return result, fmt.Errorf("invalid response origin")
		}
		destinations := map[string]bool{}
		for _, field := range group.Fields {
			if destinations[field.To] {
				return result, fmt.Errorf("%s: duplicate destination %s", group.Name, field.To)
			}
			destinations[field.To] = true
			path := field.From
			if group.Direction == "response" {
				path = field.To
			}
			resolved, err := schema.Resolve(shared.Name, path)
			if err != nil {
				return result, err
			}
			sharedField := resolved.Fields[len(resolved.Fields)-1]
			messageName := protoreflect.FullName("")
			if sharedField.Message() != nil {
				messageName = sharedField.Message().FullName()
			}
			if field.Transform == "policy_minutes" {
				if group.Direction != "request" || messageName != "kit.sandbox.v1.AutomaticAction" || field.Policy == nil || (field.Policy.Disabled != "zero" && field.Policy.Disabled != "reject") || field.Cast != "" || field.Factor != 0 {
					return result, fmt.Errorf("invalid policy_minutes mapping %s", field.From)
				}
			} else if field.Policy != nil {
				return result, fmt.Errorf("policy rules require policy_minutes transform")
			}
			if field.Transform == "duration_seconds" && (group.Direction != "request" || messageName != "google.protobuf.Duration" || field.Cast != "") {
				return result, fmt.Errorf("duration_seconds requires shared duration and no cast")
			}
			if field.Transform == "scaled_integer" && (group.Direction != "request" || sharedField.Kind() != protoreflect.DoubleKind || field.Cast != "") {
				return result, fmt.Errorf("scaled_integer requires a shared double, positive factor, and no cast")
			}
			switch field.Transform {
			case "", "whole", "divide_exactly", "multiply", "scaled_integer", "duration_seconds", "policy_minutes":
			default:
				return result, fmt.Errorf("unknown conversion %s", field.Transform)
			}
			switch field.Cast {
			case "", "integer", "number", "unsigned", "string":
			default:
				return result, fmt.Errorf("unknown cast %s", field.Cast)
			}
			if (field.Transform == "divide_exactly" || field.Transform == "multiply" || field.Transform == "scaled_integer") && field.Factor == 0 {
				return result, fmt.Errorf("%s: conversion factor must be positive", group.Name)
			}
			errorPath := field.From
			if group.Direction == "request" {
				if group.ErrorPath != "" {
					errorPath = group.ErrorPath + "." + field.From
				} else if shared.Name != "CreateOptions" {
					if root, err := schema.Message("CreateOptions"); err == nil {
						for i := 0; i < root.Fields().Len(); i++ {
							f := root.Fields().Get(i)
							if f.Message() != nil && f.Message().FullName() == resolved.Root {
								errorPath = string(f.Name()) + "." + field.From
							}
						}
					}
				}
				if group.ErrorPath != "" {
					if _, err := schema.Resolve("CreateOptions", errorPath); err != nil {
						return result, err
					}
				}
			}
			compiled.Fields = append(compiled.Fields, Mapping{Rule: field, Shared: resolved, ErrorPath: errorPath})
		}
		result.Groups = append(result.Groups, compiled)
	}
	for _, check := range rules.Checks {
		if len(check.ExclusivePolicies) > 0 || check.ForbidPolicyWith != "" {
			paths := append([]string{}, check.ExclusivePolicies...)
			if check.Path != "" {
				paths = append(paths, check.Path)
			}
			if len(check.ExclusivePolicies) > 0 && len(check.ExclusivePolicies) != 2 {
				return result, fmt.Errorf("exclusive policies requires two fields")
			}
			if check.ForbidPolicyWith != "" && check.ForbidPolicyWith != "ephemeral" && check.ForbidPolicyWith != "stopped_delete" {
				return result, fmt.Errorf("unknown policy relationship %s", check.ForbidPolicyWith)
			}
			for _, path := range paths {
				resolved, err := schema.Resolve("LifetimePolicy", path)
				if err != nil {
					return result, err
				}
				last := resolved.Fields[len(resolved.Fields)-1]
				if last.Message() == nil || last.Message().FullName() != "kit.sandbox.v1.AutomaticAction" {
					return result, fmt.Errorf("relationship requires automatic policy %s", path)
				}
			}
		} else if check.Path != "" {
			if _, err := schema.Resolve("CreateOptions", check.Path); err != nil {
				return result, err
			}
		}
	}
	return result, nil
}
