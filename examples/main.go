// This example shows the same core client API with either provider.
package main

import (
	"fmt"
	"io"
	"os"

	daytonaSDK "github.com/daytona/clients/sdk-go/pkg/daytona"
	modalSDK "github.com/modal-labs/modal-client/go"
	daytonaAdapter "github.com/sandbox-kit/kit/adapters/daytona"
	modalAdapter "github.com/sandbox-kit/kit/adapters/modal"
	"github.com/sandbox-kit/kit/core"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, output io.Writer) error {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help")) {
		_, err := fmt.Fprintln(output, "Usage: go run . <modal|daytona>\n\nInitialize your SDK, create its adapter, and use the common core client.")
		return err
	}
	if len(args) != 1 {
		return fmt.Errorf("choose one provider: modal or daytona")
	}
	switch args[0] {
	case "modal":
		return useModal(output)
	case "daytona":
		return useDaytona(output)
	default:
		return fmt.Errorf("unknown provider %q: choose modal or daytona", args[0])
	}
}

func useModal(output io.Writer) error {
	// 1. The application initializes its SDK using its own profile/credentials.
	sdkClient, err := modalSDK.NewClient()
	if err != nil {
		return fmt.Errorf("initialize Modal SDK: %w", err)
	}
	defer sdkClient.Close() // The application retains ownership of SDK resources.

	// 2. Pass the initialized SDK client to the independently installed adapter.
	provider, err := modalAdapter.New(sdkClient)
	if err != nil {
		return err
	}

	// 3. Use the same core client and public API for either provider.
	client, err := core.NewClient(provider)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "Sandbox Kit ready: %s\n", client.ProviderName())
	return err
}

func useDaytona(output io.Writer) error {
	// 1. The application initializes its SDK using its own credentials/config.
	sdkClient, err := daytonaSDK.NewClient()
	if err != nil {
		return fmt.Errorf("initialize Daytona SDK: %w", err)
	}

	// 2. The adapter borrows the configured SDK client, preserving authentication.
	provider, err := daytonaAdapter.New(sdkClient)
	if err != nil {
		return err
	}

	// 3. Only the adapter changes; the core client API remains the same.
	client, err := core.NewClient(provider)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(output, "Sandbox Kit ready: %s\n", client.ProviderName())
	return err
}
