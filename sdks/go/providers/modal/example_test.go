package modal_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func ExampleNew() {
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
		Lifetime: &sandbox.LifetimePolicy{
			MaximumLifetime: sandbox.Duration(5 * time.Minute),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(instance.ID())
}

func ExampleNew_attach() {
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

func ExampleNew_terminate() {
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
	_, err = instance.Terminate(context.Background(), &sandbox.TerminateOptions{
		Wait: sandbox.Value(true),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleNew_errors() {
	_, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Endpoint: sandbox.Value("https://example.com"),
	})
	var detail *sandbox.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Field)
	}
	// Output:
	// endpoint
}
