package modal

import (
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestCloneMessageKeepsExplicitEmptyCollections(t *testing.T) {
	src := &sandbox.CreateOptions{
		Environment: map[string]string{},
		Runtime:     &sandbox.RuntimeConfig{Entrypoint: []string{}},
	}
	dst := cloneMessage(src)
	if dst == nil || dst == src || dst.Environment == nil || dst.Runtime == nil || dst.Runtime.Entrypoint == nil {
		t.Fatalf("empty collections dropped: %#v", dst)
	}
	dst.Environment["MODE"] = "changed"
	if len(src.Environment) != 0 {
		t.Fatal("clone shares the source map")
	}
	if cloneMessage((*sandbox.CreateOptions)(nil)) != nil {
		t.Fatal("nil source returned a message")
	}
}
