package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() (result error) {
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}
	client, err := sandbox.NewClient(sandbox.Config{Provider: modal.New(), Auth: &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: os.Getenv("MODAL_TOKEN_ID"), Secret: os.Getenv("MODAL_TOKEN_SECRET")}}, Scope: &sandbox.Scope{AppName: sandbox.Value(os.Getenv("MODAL_APP"))}, Timeout: sandbox.Value(2 * time.Minute)})
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(context.Background()); result == nil {
			result = err
		}
	}()
	instance, err := client.Attach(context.Background(), os.Getenv("MODAL_SANDBOX_ID"))
	if err != nil {
		return err
	}
	fmt.Printf("Modal handle attached: id=%s capabilities=%v\n", instance.ID(), instance.Capabilities())
	return nil
}
