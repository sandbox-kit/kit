package providergen

import (
	"fmt"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateErrors emits native provider-error classification from the compiled
// YAML rules. The original native error remains the SandboxError cause.
func GenerateErrors(plugin *protogen.Plugin, file *protogen.File, provider model.Provider) error {
	binding := provider.Rules.Errors["typescript"]
	if len(provider.Errors.HTTP) > 0 && binding.Target.Import == "" {
		return fmt.Errorf("provider %s: HTTP error rules require a TypeScript error target", provider.Rules.Provider)
	}

	g := declarationgen.GeneratedFile(plugin, "provider.errors.gen.ts", string(file.Desc.Path())+" + specs/providers/"+provider.Rules.Provider+".yaml")
	imports := map[string][]string{}
	for _, native := range binding.Types {
		if strings.TrimSpace(native.Import) == "" || !token.IsIdentifier(native.Name) {
			return fmt.Errorf("provider %s: invalid TypeScript error type", provider.Rules.Provider)
		}
		imports[native.Import] = append(imports[native.Import], native.Name)
	}
	if binding.Target.Import != "" {
		if !token.IsIdentifier(binding.Target.Name) {
			return fmt.Errorf("provider %s: invalid TypeScript error target", provider.Rules.Provider)
		}
		imports[binding.Target.Import] = append(imports[binding.Target.Import], binding.Target.Name)
	}
	paths := make([]string, 0, len(imports))
	for path := range imports {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		names := uniqueSorted(imports[path])
		g.P("import { ", strings.Join(names, ", "), " } from ", strconv.Quote(path), ";")
	}
	g.P(`import { ErrorKind, SandboxError, sandboxError, type ErrorInfo } from "@sandbox-kit/sandbox";`)
	g.P()
	functionName := "map" + naming.Pascal(provider.Rules.Provider) + "Error"
	g.P("export function ", functionName, "(error: unknown, operation: string): SandboxError {")
	g.P("  if (error instanceof SandboxError) return error;")
	g.P(`  if (error instanceof DOMException && error.name === "AbortError") return create(ErrorKind.Canceled, error, operation);`)
	g.P(`  if (error instanceof DOMException && error.name === "TimeoutError") return create(ErrorKind.Timeout, error, operation);`)
	g.P(`  let kind: ErrorInfo["kind"] = ErrorKind.Unknown;`)
	for _, native := range binding.Types {
		classification, ok := nativeKind(provider, native.Rule)
		if !ok {
			return fmt.Errorf("provider %s: missing compiled error rule %s", provider.Rules.Provider, native.Rule)
		}
		g.P("  if (kind === ErrorKind.Unknown && error instanceof ", native.Name, ") kind = ErrorKind.", kindName(classification), ";")
	}
	g.P("  let statusCode: number | undefined;")
	g.P("  let providerCode: string | undefined;")
	g.P("  let providerSource: string | undefined;")
	if binding.Target.Import != "" {
		for _, field := range []string{"status", "code", "source"} {
			if !token.IsIdentifier(binding.Target.Fields[field]) {
				return fmt.Errorf("provider %s: missing TypeScript error member %s", provider.Rules.Provider, field)
			}
		}
		status := binding.Target.Fields["status"]
		code := binding.Target.Fields["code"]
		source := binding.Target.Fields["source"]
		g.P("  if (error instanceof ", binding.Target.Name, ") {")
		g.P("    statusCode = error.", status, "; providerCode = error.", code, "; providerSource = error.", source, ";")
		if len(provider.Errors.Codes) > 0 {
			g.P("    if (kind === ErrorKind.Unknown) switch (error.", code, ") {")
			for _, rule := range provider.Errors.Codes {
				g.P("      case ", strconv.Quote(rule.CodeText), ": kind = ErrorKind.", kindName(rule.Kind), "; break;")
			}
			g.P("    }")
		}
		if len(provider.Errors.HTTP) > 0 {
			g.P("    if (kind === ErrorKind.Unknown) switch (error.", status, ") {")
			for _, rule := range provider.Errors.HTTP {
				g.P("      case ", rule.Code, ": kind = ErrorKind.", kindName(rule.Kind), "; break;")
			}
			g.P("    }")
		}
		if provider.Errors.Default != nil {
			g.P("    if (kind === ErrorKind.Unknown && error.", status, " === 0 && !error.", code, ") kind = ErrorKind.", kindName(*provider.Errors.Default), ";")
		}
		g.P("  }")
	}
	g.P("  return create(kind, error, operation, statusCode, providerCode, providerSource);")
	g.P("}")
	g.P()
	g.P("function create(kind: ErrorInfo[\"kind\"], error: unknown, operation: string, statusCode?: number, providerCode?: string, providerSource?: string): SandboxError {")
	g.P("  return sandboxError({ kind, provider: ", strconv.Quote(provider.Rules.Provider), ", operation, field: \"\", message: error instanceof Error ? error.message : ", strconv.Quote("sandbox-kit "+provider.Rules.Provider+": provider operation failed"), ", ...(statusCode === undefined ? {} : { statusCode }), ...(providerCode === undefined ? {} : { providerCode }), ...(providerSource === undefined ? {} : { providerSource }) }, error);")
	g.P("}")
	return nil
}

func nativeKind(provider model.Provider, rule string) (model.Kind, bool) {
	for _, candidate := range provider.Errors.Native {
		if candidate.Rule == rule {
			return candidate.Kind, true
		}
	}
	return model.Kind{}, false
}

func kindName(kind model.Kind) string {
	name := string(kind.Name)
	if index := strings.LastIndexByte(name, '.'); index >= 0 {
		name = name[index+1:]
	}
	name = strings.TrimPrefix(name, "ERROR_KIND_")
	return naming.Pascal(strings.ToLower(name))
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}
