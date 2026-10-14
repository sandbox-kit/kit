package providergen

import (
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestKindName(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string
		want string
	}{
		{name: "kit.sandbox.v1.ERROR_KIND_INVALID_ARGUMENT", want: "InvalidArgument"},
		{name: "kit.sandbox.v1.ERROR_KIND_RATE_LIMITED", want: "RateLimited"},
	} {
		if got := kindName(model.Kind{Name: protoreflectName(test.name)}); got != test.want {
			t.Fatalf("kindName(%q) = %q, want %q", test.name, got, test.want)
		}
	}
}

func protoreflectName(value string) protoreflect.FullName { return protoreflect.FullName(value) }
