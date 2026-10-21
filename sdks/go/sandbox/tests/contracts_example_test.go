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

func ExampleCloneSandboxInfo() {
	metadata, err := sandbox.MetadataObjectFromMap(map[string]any{"generation": int64(9007199254740993)})
	if err != nil {
		panic(err)
	}
	original := &sandbox.SandboxInfo{Id: "example", ProviderMetadata: metadata, Origins: map[string]sandbox.ValueOrigin{"id": sandbox.ValueOriginProvider}}
	copied := sandbox.CloneSandboxInfo(original)
	original.ProviderMetadata.Values["generation"] = &sandbox.MetadataValue{Value: &sandbox.MetadataSignedInteger{SignedInteger: 0}}
	copiedMap, err := sandbox.MetadataObjectAsMap(copied.ProviderMetadata)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T %v %v\n", copiedMap["generation"], copiedMap["generation"], copied.Origins["id"])
	// Output: int64 9007199254740993 VALUE_ORIGIN_PROVIDER
}
