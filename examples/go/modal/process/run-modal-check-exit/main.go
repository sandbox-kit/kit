// Check a non-zero exit code from a Modal command.
package main

import (
	"context"
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
	instance, err := client.Attach(context.Background(), os.Getenv("MODAL_SANDBOX_ID"))
	if err != nil {
		log.Fatal(err)
	}
	result, err := instance.Process().Run(context.Background(), &sandbox.ExecOptions{
		Shell: sandbox.Value("exit 7"),
	})
	if err != nil {
		log.Fatal(err)
	}
	if result.GetExitCode() != 7 {
		log.Fatalf("want exit 7, got %d", result.GetExitCode())
	}
	fmt.Printf("command finished with exit=%d\n", result.GetExitCode())
}
