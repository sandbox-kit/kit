package typesgen

import (
	"strings"

	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	"google.golang.org/protobuf/compiler/protogen"
)

// commentBefore writes a wrapped godoc comment. The Go name starts the first line.
func commentBefore(g *protogen.GeneratedFile, comments protogen.Comments, name string) {
	text := strings.TrimSpace(string(comments))
	if text == "" {
		return
	}
	lines := strings.Split(text, "\n")
	var b strings.Builder
	for i, line := range lines {
		line = strings.TrimRight(line, " \t")
		if strings.HasPrefix(line, " ") {
			line = line[1:]
		}
		if i == 0 {
			line = strings.TrimSpace(line)
			if name != "" && !strings.HasPrefix(line, name+" ") {
				line = name + " " + line
			}
		}
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
	}
	declarationgen.WriteComment(g, "", b.String())
}
