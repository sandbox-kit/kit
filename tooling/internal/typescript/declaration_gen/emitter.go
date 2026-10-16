package declarationgen

import (
	"fmt"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

// Emitter writes readable TypeScript source blocks. It keeps formatting and
// template execution in one place so semantic generators do not have to build
// complete source lines through repeated GeneratedFile.P calls.
type Emitter struct {
	file *protogen.GeneratedFile
}

func NewEmitter(file *protogen.GeneratedFile) *Emitter {
	return &Emitter{file: file}
}

func (e *Emitter) Line(parts ...any) {
	e.file.P(parts...)
}

func (e *Emitter) Blank() {
	e.file.P()
}

// Block emits a dedented source block and preserves intentional blank lines.
func (e *Emitter) Block(source string) {
	e.IndentedBlock(source, "")
}

// IndentedBlock emits a dedented source block below an existing TypeScript
// scope, adding prefix to every non-empty line.
func (e *Emitter) IndentedBlock(source, prefix string) {
	source = strings.Trim(source, "\n")
	if source == "" {
		return
	}
	lines := strings.Split(source, "\n")
	indent := commonIndent(lines)
	for _, line := range lines {
		if indent > 0 && len(line) >= indent {
			line = line[indent:]
		}
		line = strings.TrimRight(line, " \t")
		if line != "" {
			line = prefix + line
		}
		e.file.P(line)
	}
}

// Template renders a TypeScript source block from semantic generator data.
func (e *Emitter) Template(name, source string, data any) error {
	return e.IndentedTemplate(name, source, data, "")
}

// IndentedTemplate renders semantic data and emits it below an existing scope.
func (e *Emitter) IndentedTemplate(name, source string, data any, prefix string) error {
	tmpl, err := template.New(name).Option("missingkey=error").Parse(source)
	if err != nil {
		return fmt.Errorf("parse TypeScript template %s: %w", name, err)
	}
	var rendered strings.Builder
	if err := tmpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("render TypeScript template %s: %w", name, err)
	}
	e.IndentedBlock(rendered.String(), prefix)
	return nil
}

func commonIndent(lines []string) int {
	indent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		width := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent < 0 || width < indent {
			indent = width
		}
	}
	if indent < 0 {
		return 0
	}
	return indent
}
