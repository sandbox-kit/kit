package main

import (
	"context"
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
	instance, err := client.Get(context.Background(), os.Getenv("BL_SANDBOX_ID"))
	if err != nil {
		log.Fatal(err)
	}
	cwd := os.Getenv("BL_WORKING_DIRECTORY")
	if cwd == "" {
		cwd = "/tmp"
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Argv:             []string{"pwd"},
		WorkingDirectory: sandbox.Value(cwd),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("exit=%d cwd_out=%q\n", result.GetExitCode(), string(result.GetStdout()))
}
