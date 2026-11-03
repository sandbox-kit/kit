package sandbox_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleDuration() {
	fmt.Println(sandbox.Duration(2 * time.Minute).AsDuration())
	// Output:
	// 2m0s
}

func ExampleValue() {
	resources := &sandbox.Resources{
		CpuCores:  sandbox.Value(2.0),
		MemoryMib: sandbox.Value(uint64(4096)),
	}
	fmt.Println(*resources.CpuCores, *resources.MemoryMib)
	// Output:
	// 2 4096
}

func ExampleValidateConfig() {
	err := sandbox.ValidateConfig(&sandbox.Config{
		Endpoint: sandbox.Value("http://example.com"),
	})
	fmt.Println(err != nil)
	// Output:
	// true
}

func ExampleValidateCreateOptions() {
	err := sandbox.ValidateCreateOptions(&sandbox.CreateOptions{
		Resources: &sandbox.Resources{CpuCores: sandbox.Value(0.0)},
	})
	fmt.Println(err != nil)
	// Output:
	// true
}

func ExampleNewError() {
	err := sandbox.NewError(&sandbox.ErrorInfo{
		Kind:    sandbox.ErrorKindUnsupported,
		Field:   sandbox.CreateFieldResourcesDiskMib,
		Message: "disk overrides are unsupported",
	}, nil)
	var detail *sandbox.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Field)
	}
	// Output:
	// resources.disk_mib
}

func ExampleMetadataObjectAsMap() {
	object := &sandbox.MetadataObject{Values: map[string]*sandbox.MetadataValue{
		"team": {Value: &sandbox.MetadataText{Text: "platform"}},
	}}
	values, err := sandbox.MetadataObjectAsMap(object)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(values["team"])
	// Output:
	// platform
}

func ExampleExecOptions() {
	opts := &sandbox.ExecOptions{
		Argv:    []string{"echo", "hello from kit"},
		Timeout: sandbox.Duration(5 * time.Second),
	}
	fmt.Println(opts.Argv[0], opts.Timeout.AsDuration())
	// Output:
	// echo 5s
}

func ExampleProcess_Run() {
	// Obtain Process from a live Sandbox handle (Create, Get, or Attach).
	// A nil handle fails before any provider call.
	_, err := (*sandbox.Sandbox)(nil).Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv: []string{"echo", "hello from kit"},
	})
	var detail *sandbox.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Kind)
	}
	// Output:
	// ERROR_KIND_FAILED_PRECONDITION
}

func ExampleFilesystem_WriteFile() {
	// Obtain Filesystem from a live Sandbox handle (Create, Get, or Attach).
	err := (*sandbox.Sandbox)(nil).Filesystem().WriteFile(context.Background(), "/tmp/note.txt", []byte("hello\n"), nil)
	var detail *sandbox.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Kind)
	}
	// Output:
	// ERROR_KIND_FAILED_PRECONDITION
}
