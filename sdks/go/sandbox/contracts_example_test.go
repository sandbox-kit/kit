package sandbox_test

import (
	"errors"
	"fmt"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNewError() {
	native := errors.New("original provider failure")
	err := sandbox.NewError(sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindUnsupported,
		Provider:  "modal",
		Operation: "create",
		Field:     "resources.disk_mib",
		Message:   "disk overrides are unavailable",
	}, native)
	var detail *sandbox.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Kind, detail.Field, errors.Is(err, native))
	}
	// Output: ERROR_KIND_UNSUPPORTED resources.disk_mib true
}

func ExampleSandboxInfo_Clone() {
	original := &sandbox.SandboxInfo{ID: "example", ProviderMetadata: map[string]any{"generation": int64(9007199254740993)}, Origins: map[string]sandbox.ValueOrigin{"id": sandbox.ValueOriginProvider}}
	copied, err := original.Clone()
	if err != nil {
		panic(err)
	}
	original.ProviderMetadata["generation"] = int64(0)
	fmt.Printf("%T %v %v\n", copied.ProviderMetadata["generation"], copied.ProviderMetadata["generation"], copied.Origins["id"])
	// Output: int64 9007199254740993 VALUE_ORIGIN_PROVIDER
}
