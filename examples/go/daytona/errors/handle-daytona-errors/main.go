// Demonstrate common error handling with the Daytona integration.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("copy .env.example to .env and fill in its values: %v", err)
	}
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.APIKey{
				ApiKey: &sandbox.APIKeyCredentials{
					Key: os.Getenv("DAYTONA_API_KEY"),
				},
			},
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	}
	if endpoint := os.Getenv("DAYTONA_API_URL"); endpoint != "" {
		config.Endpoint = sandbox.Value(endpoint)
	}
	if target := os.Getenv("DAYTONA_TARGET"); target != "" {
		config.Region = sandbox.Value(target)
	}
	client, err := sandbox.NewClient(daytona.New(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	// Invalid CPU demonstrates how an application handles configuration errors.
	_, err = client.Create(context.Background(), &sandbox.CreateOptions{
		Resources: &sandbox.Resources{CpuCores: sandbox.Value(-1.0)},
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
