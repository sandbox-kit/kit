package sandbox_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNewClient() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
		}},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNewClient_modal() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{
			TokenPair: &sandbox.TokenPairCredentials{
				Id:     os.Getenv("MODAL_TOKEN_ID"),
				Secret: os.Getenv("MODAL_TOKEN_SECRET"),
			},
		}},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value("main"),
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
		}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleClient_Create() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{Kind: &sandbox.Image{
			Image: &sandbox.ImageSource{Reference: "python:3.11"},
		}},
		Resources: &sandbox.Resources{
			CpuCores:  sandbox.Value(2.0),
			MemoryMib: sandbox.Value(uint64(4096)),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleClient_Get() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleClient_Attach() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleClient_Close() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Close releases the local SDK. The cloud sandbox keeps running.
	err = client.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Info() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	info := instance.Info()
	fmt.Println(info.GetId(), instance.ProviderName())
}

func ExampleSandbox_Supports() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	if instance.Supports(sandbox.LifecycleCapabilityPause) {
		err = instance.Pause(context.Background(), nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ExampleSandbox_Refresh() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Refresh(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Start() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Start(context.Background(), &sandbox.StartOptions{
		Timeout: sandbox.Duration(time.Minute),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Stop() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Stop(context.Background(), &sandbox.StopOptions{
		Force: sandbox.Value(true),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Pause() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Pause(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Archive() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Archive(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Delete() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Get(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = instance.Delete(context.Background(), &sandbox.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleSandbox_Terminate() {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP_NAME"))},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Attach(context.Background(), os.Getenv("SANDBOX_ID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := instance.Terminate(context.Background(), &sandbox.TerminateOptions{
		Wait: sandbox.Value(true),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result.GetExitCode())
}

func ExampleSandbox_Close() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close(context.Background())

	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Close drops the local handle. Daytona keeps the sandbox.
	err = instance.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

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

func ExampleMetadataObjectFromMap() {
	object, err := sandbox.MetadataObjectFromMap(map[string]any{"team": "platform"})
	if err != nil {
		fmt.Println(err)
		return
	}
	values, err := sandbox.MetadataObjectAsMap(object)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(values["team"])
	// Output:
	// platform
}
