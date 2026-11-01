package modal

import (
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestMakeDirectoryParams(t *testing.T) {
	for _, options := range []*sandbox.MakeDirectoryOptions{
		nil,
		{},
		{Parents: sandbox.Value(true)},
	} {
		params := makeDirectoryParams(options)
		if params == nil || params.CreateParents == nil || !*params.CreateParents {
			t.Fatalf("parents true: %#v from %#v", params, options)
		}
	}
	params := makeDirectoryParams(&sandbox.MakeDirectoryOptions{Parents: sandbox.Value(false)})
	if params == nil || params.CreateParents == nil || *params.CreateParents {
		t.Fatalf("parents false: %#v", params)
	}
}
