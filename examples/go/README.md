# Go example: common client configuration and sandbox creation

One runnable program demonstrates the current API. It imports only sandbox and the
optional Sandbox Kit adapters; adapters initialize their official SDKs internally.
The previous borrowed-SDK constructor flow has been removed.

## Run locally without credentials

```sh
cd examples/go
go run . --help
go run . demo
GOWORK=off go test ./...
```

The demo implements the same provider/backend contracts but creates no cloud
resources. Expected output: `Simulated sandbox created: demo-sandbox (demo)`.

## Common client config

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: daytona.New(),
    Auth: &sandbox.AuthConfig{
        APIKey: &sandbox.APIKeyCredentials{Key: os.Getenv("DAYTONA_API_KEY")},
    },
    Endpoint: sandbox.Value("https://app.daytona.io/api"),
    Region: sandbox.Value("us"),
    Timeout: sandbox.Value(30 * time.Second),
})
if err != nil {
    return err
}
defer client.Close(context.Background())
```

Use a target/endpoint valid for your account; the strings above illustrate fields.
`Config` is generated from protobuf and YAML. `Provider` is a typed factory
attachment, not an arbitrary value or an initialized SDK client. A name-only
interface does not satisfy the initialization contract. Valid custom factories
can implement `sandbox.Provider` and return the full `sandbox.Backend` contract.

For Modal, change the selected provider and supported authentication:

```go
client, err := sandbox.NewClient(sandbox.Config{
    Provider: modal.New(),
    Auth: &sandbox.AuthConfig{
        TokenPair: &sandbox.TokenPairCredentials{
            ID: os.Getenv("MODAL_TOKEN_ID"),
            Secret: os.Getenv("MODAL_TOKEN_SECRET"),
        },
    },
    Scope: &sandbox.Scope{
        AppName: sandbox.Value("my-existing-app"),
        Environment: sandbox.Value("main"),
    },
    Timeout: sandbox.Value(30 * time.Second),
})
```

Check `err` and close this client when finished too. App context belongs in
`Config`, not in each creation request. Modal uses token-pair or OAuth
refresh credentials; Daytona uses API key or bearer/JWT credentials. Daytona
bearer authentication requires an organization in `Scope.OrganizationID` or
its SDK environment settings. Omit `Auth` entirely to use SDK profile/environment
resolution; an explicit empty `AuthConfig` is invalid.

Modal does not expose a per-client endpoint override in the pinned Go SDK.
Setting `Config.Endpoint` for Modal returns a local error before SDK
initialization. Daytona supports endpoint and region/target directly. Modal
region is a default placement for later creation requests, which may override it.
Unsupported authentication/context settings are rejected before SDK initialization.

## Initialize a real provider

Configure credentials in your environment, then:

```sh
go run . modal --app my-existing-app
go run . daytona
```

These initialize SDKs and print `Sandbox Kit ready: <provider>`. They do not
verify credentials through a cloud API or provision a sandbox.

## Explicit cloud creation

```sh
go run . modal --app my-existing-app --create --image python:3.11
go run . daytona --create --image python:3.11
```

`--create` provisions real resources. Additional flags include `--region`,
`--endpoint` (Daytona), `--environment` (Modal), `--organization` (Daytona), and
`--timeout`. SDK environment/profile defaults remain available for omitted settings.
`main.go` checks construction and creation errors and propagates SDK cleanup errors.

Both providers use the same creation API:

```go
instance, err := client.Create(ctx, &sandbox.CreateOptions{
    Source: &sandbox.SandboxSource{
        Image: &sandbox.ImageSource{Reference: "python:3.11"},
    },
    Resources: &sandbox.Resources{
        CPUCores: sandbox.Value(2.0),
        MemoryMiB: sandbox.Value(uint64(4096)),
    },
    Environment: map[string]string{"MODE": "development"},
})
if err != nil {
    return err
}
fmt.Println(instance.ID(), instance.ProviderName())
```

An explicit `Provisioning.Timeout` overrides the client default, including zero.
The caller's context deadline still applies. `Close` releases SDK resources; it
does not delete the created sandbox. Sandbox lifecycle methods are a later increment.

See [client configuration](../../docs/configuration.md),
[creation support](../../docs/sandbox-creation.md), and
[generator inputs](../../docs/code-generation.md).
The example uses local module replacements; packages are not published yet.
