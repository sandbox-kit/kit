package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/blaxel"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("copy .env.example to .env and fill in its values: %v", err)
	}
	region := os.Getenv("BL_REGION")
	if region == "" {
		region = "us-pdx-1"
	}
	config := &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.APIKey{
				ApiKey: &sandbox.APIKeyCredentials{
					Key: os.Getenv("BL_API_KEY"),
				},
			},
		},
		Scope: &sandbox.Scope{
			Workspace: sandbox.Value(os.Getenv("BL_WORKSPACE")),
		},
		Region:  sandbox.Value(region),
		Timeout: sandbox.Duration(2 * time.Minute),
	}
	client, err := sandbox.NewClient(blaxel.New(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close(context.Background())
	_, err = client.Create(context.Background(), &sandbox.CreateOptions{
		Name:      sandbox.Value("kit-handle-blaxel-errors"),
		Source:    &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "blaxel/base-image:latest"}}},
		Resources: &sandbox.Resources{CpuCores: sandbox.Value(2.0)},
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
		default:
			fmt.Fprintln(os.Stdout, err)
		}
		return
	}
	fmt.Fprintln(os.Stdout, "Sandbox created.")
}
