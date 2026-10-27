package daytona_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNew() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{Credentials: &sandbox.APIKey{
			ApiKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
		}},
		Region:  sandbox.Value("us"),
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
	fmt.Println(instance.ProviderName(), instance.ID())
}

func ExampleNew_image() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
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
			DiskMib:   sandbox.Value(uint64(8192)),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNew_get() {
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
	if err = instance.Refresh(context.Background()); err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_stop() {
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
	err = instance.Stop(context.Background(), &sandbox.StopOptions{Force: sandbox.Value(false)})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_info() {
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

func ExampleNew_supports() {
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

func ExampleNew_start() {
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

func ExampleNew_pause() {
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

func ExampleNew_archive() {
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

func ExampleNew_delete() {
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

func ExampleNew_close() {
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
	// Close returns without a provider request. The Daytona sandbox keeps running.
	err = instance.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_errors() {
	client, err := sandbox.NewClient(daytona.New(), &sandbox.Config{})
	if err != nil {
		var detail *sandbox.Error
		if errors.As(err, &detail) {
			fmt.Println(detail.Kind, detail.Field)
		}
		return
	}
	defer client.Close(context.Background())
}
