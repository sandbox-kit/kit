package sandbox

import (
	"math"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestCopiedMetadataRetainsNumbersAndOwnership(t *testing.T) {
	big := uint64(math.MaxUint64)
	metadata, err := MetadataObjectFromMap(map[string]any{"large": big, "signed": int64(9007199254740993), "nested": []any{map[string]any{"bytes": []byte{1, 2}}}})
	if err != nil {
		t.Fatal(err)
	}
	source := &SandboxInfo{Id: "id", ProviderMetadata: metadata, Origins: map[string]ValueOrigin{"id": ValueOriginProvider}}
	copied := cloneSandboxInfo(source)
	if !proto.Equal(source, copied) {
		t.Fatal("types or values changed")
	}
	sourceMap, err := MetadataObjectAsMap(source.ProviderMetadata)
	if err != nil {
		t.Fatal(err)
	}
	sourceMap["large"] = uint64(0)
	sourceMap["nested"].([]any)[0].(map[string]any)["bytes"].([]byte)[0] = 9
	source.Origins["id"] = ValueOriginRequest
	copiedMap, err := MetadataObjectAsMap(copied.ProviderMetadata)
	if err != nil {
		t.Fatal(err)
	}
	if copiedMap["large"] != big || copiedMap["nested"].([]any)[0].(map[string]any)["bytes"].([]byte)[0] != 1 || copied.Origins["id"] != ValueOriginProvider {
		t.Fatal("copy aliases original")
	}
}
