// Demonstrate common error handling with the Modal integration.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("copy .env.example to .env and fill in its values: %v", err)
	}
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.TokenPair{
				TokenPair: &sandbox.TokenPairCredentials{
					Id:     os.Getenv("MODAL_TOKEN_ID"),
					Secret: os.Getenv("MODAL_TOKEN_SECRET"),
				},
			},
		},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value(environment),
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	// Modal does not support a disk allocation override at creation.
	_, err = client.Create(context.Background(), &sandbox.CreateOptions{
		Source:    &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}}},
		Resources: &sandbox.Resources{DiskMib: sandbox.Value(uint64(1024))},
	})
	if err != nil {
		var sandboxError *sandbox.Error
		if !errors.As(err, &sandboxError) {
			fmt.Fprintln(os.Stdout, err)
			return
		}
		switch sandboxError.Kind {
		case sandbox.ErrorKindInvalidArgument, sandbox.ErrorKindUnsupported:
			fmt.Fprintf(os.Stdout, "Check %s: %s\n", sandboxError.Field, sandboxError.Message)
		case sandbox.ErrorKindAuthentication:
			fmt.Fprintln(os.Stdout, "Check your provider credentials.")
		case sandbox.ErrorKindPermissionDenied:
			fmt.Fprintln(os.Stdout, "Your account does not have permission for this operation.")
		case sandbox.ErrorKindNotFound:
			fmt.Fprintln(os.Stdout, "The requested provider resource was not found.")
		case sandbox.ErrorKindRateLimited, sandbox.ErrorKindResourceExhausted:
			fmt.Fprintln(os.Stdout, "The provider's rate or resource limit was reached.")
		case sandbox.ErrorKindCanceled:
			fmt.Fprintln(os.Stdout, "The operation was canceled.")
		case sandbox.ErrorKindTimeout:
			fmt.Fprintln(os.Stdout, "The operation timed out.")
		default:
			fmt.Fprintln(os.Stdout, err)
		}
		return
	}
	fmt.Fprintln(os.Stdout, "Sandbox created.")
}
