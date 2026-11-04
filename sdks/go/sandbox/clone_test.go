package sandbox

import (
	"math"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestCloneExecOptionsPreservesEmptyEnvironment(t *testing.T) {
	options := &ExecOptions{Argv: []string{"echo"}, Environment: map[string]string{}}
	cloned := cloneExecOptions(options)
	if cloned.Environment == nil {
		t.Fatal("explicit empty environment became nil")
	}
	options.Environment["X"] = "1"
	if _, ok := cloned.Environment["X"]; ok {
		t.Fatal("clone aliases original environment")
	}
}

func TestCloneExecResultPreservesEmptyOutput(t *testing.T) {
	result := &ExecResult{ExitCode: 0, Stdout: []byte{}, Stderr: []byte{}}
	cloned := cloneExecResult(result)
	if cloned.Stdout == nil || cloned.Stderr == nil {
		t.Fatalf("empty output became nil: stdout=%#v stderr=%#v", cloned.Stdout, cloned.Stderr)
	}
}

func TestCopiedMetadataRetainsNumbersAndOwnership(t *testing.T) {
	big := uint64(math.MaxUint64)
	metadata := &MetadataObject{Values: map[string]*MetadataValue{
		"large":  {Value: &MetadataUnsignedInteger{UnsignedInteger: big}},
		"signed": {Value: &MetadataSignedInteger{SignedInteger: 9007199254740993}},
		"nested": {Value: &MetadataListValue{List: &MetadataList{Values: []*MetadataValue{
			{Value: &MetadataObjectValue{Object: &MetadataObject{Values: map[string]*MetadataValue{
				"bytes": {Value: &MetadataBinary{Binary: []byte{1, 2}}},
			}}}},
		}}}},
	}}
	source := &SandboxInfo{Id: "id", ProviderMetadata: metadata}
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
	copiedMap, err := MetadataObjectAsMap(copied.ProviderMetadata)
	if err != nil {
		t.Fatal(err)
	}
	if copiedMap["large"] != big || copiedMap["nested"].([]any)[0].(map[string]any)["bytes"].([]byte)[0] != 1 {
		t.Fatal("copy aliases original")
	}
}
