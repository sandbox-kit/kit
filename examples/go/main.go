// One example: common client configuration, optional providers, and sandbox creation.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/daytona"
	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type settings struct {
	create                                                  bool
	image, app, environment, organization, region, endpoint string
	timeout                                                 time.Duration
}

func run(args []string, output io.Writer) (result error) {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help")) {
		_, err := fmt.Fprintln(output, "Usage: go run . <modal|daytona|demo> [flags]\n\nCommon client config initializes the provider SDK. Add --create to provision; demo simulates creation.\nFlags: --create --image --app --environment --organization --region --endpoint --timeout")
		return err
	}
	name := args[0]
	if name != "modal" && name != "daytona" && name != "demo" {
		return fmt.Errorf("unknown provider %q: choose modal, daytona, or demo", name)
	}
	var options settings
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(output)
	flags.BoolVar(&options.create, "create", false, "Create a real sandbox with the selected provider")
	flags.StringVar(&options.image, "image", "python:3.11", "Registry image reference")
	flags.StringVar(&options.app, "app", os.Getenv("SANDBOX_APP_NAME"), "Existing application name (Modal creation)")
	flags.StringVar(&options.environment, "environment", "", "Provider environment")
	flags.StringVar(&options.organization, "organization", os.Getenv("DAYTONA_ORGANIZATION_ID"), "Organization scope (Daytona bearer authentication)")
	flags.StringVar(&options.region, "region", "", "Provider region/target")
	flags.StringVar(&options.endpoint, "endpoint", "", "Provider API endpoint (Daytona)")
	flags.DurationVar(&options.timeout, "timeout", 30*time.Second, "Default Kit operation timeout; zero disables the additional deadline")
	if err := flags.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v", flags.Args())
	}
	config := clientConfig(name, options)
	client, err := sandbox.NewClient(config)
	if err != nil {
		return err
	}
	defer func() { result = errors.Join(result, client.Close(context.Background())) }()
	if name != "demo" && !options.create {
		_, err := fmt.Fprintf(output, "Sandbox Kit ready: %s\n", client.ProviderName())
		return err
	}
	instance, err := client.Create(context.Background(), &sandbox.CreateOptions{
		Name: sandbox.Value("agent-workspace"), Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: options.image}},
		Resources:   &sandbox.Resources{CPUCores: sandbox.Value(2.0), MemoryMiB: sandbox.Value(uint64(4096))},
		Environment: map[string]string{"MODE": "development"}, Labels: map[string]string{"team": "kit"},
	})
	if err != nil {
		return err
	}
	label := "Sandbox created"
	if name == "demo" {
		label = "Simulated sandbox created"
	}
	_, err = fmt.Fprintf(output, "%s: %s (%s)\n", label, instance.ID(), instance.ProviderName())
	return err
}
func clientConfig(name string, options settings) sandbox.Config {
	config := sandbox.Config{Timeout: sandbox.Value(options.timeout)}
	if options.region != "" {
		config.Region = sandbox.Value(options.region)
	}
	if options.endpoint != "" {
		config.Endpoint = sandbox.Value(options.endpoint)
	}
	scope := &sandbox.Scope{}
	if options.app != "" {
		scope.AppName = sandbox.Value(options.app)
	}
	if options.environment != "" {
		scope.Environment = sandbox.Value(options.environment)
	}
	if options.organization != "" {
		scope.OrganizationID = sandbox.Value(options.organization)
	}
	if scope.AppName != nil || scope.Environment != nil || scope.OrganizationID != nil {
		config.Scope = scope
	}
	switch name {
	case "modal":
		config.Provider = modal.New()
		id, secret := os.Getenv("MODAL_TOKEN_ID"), os.Getenv("MODAL_TOKEN_SECRET")
		if id != "" || secret != "" {
			config.Auth = &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: id, Secret: secret}}
		}
	case "daytona":
		config.Provider = daytona.New()
		if key := os.Getenv("DAYTONA_API_KEY"); key != "" {
			config.Auth = &sandbox.AuthConfig{APIKey: &sandbox.APIKeyCredentials{Key: key}}
		} else if token := os.Getenv("DAYTONA_JWT_TOKEN"); token != "" {
			config.Auth = &sandbox.AuthConfig{BearerToken: &sandbox.BearerTokenCredentials{Token: token}}
		}
	case "demo":
		config.Provider = &demoProvider{}
	}
	return config
}

// A simulated provider implementing the same initialization and operation contracts.
// It allocates no cloud resources and needs no credentials.
type demoProvider struct{}

func (*demoProvider) Name() string { return "demo" }
func (*demoProvider) NewClient(*sandbox.Config) (sandbox.Backend, error) {
	return &demoBackend{}, nil
}

type demoBackend struct{}

func (*demoBackend) Name() string { return "demo" }
func (*demoBackend) Create(_ context.Context, r *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "demo-sandbox", Provider: "demo", Name: r.Name, Labels: r.Labels}}, nil
}
func (*demoBackend) Close(context.Context) error { return nil }
