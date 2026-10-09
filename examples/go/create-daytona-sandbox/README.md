# Create a Daytona sandbox

Authenticate with Daytona, create one sandbox from its default snapshot, and
print its ID. This standalone Go module imports Sandbox Kit and only the
Daytona provider.

## Configure credentials

Requires Go 1.26.1 or newer and this checkout. From the repository root:

```sh
cd examples/go/create-daytona-sandbox
cp .env.example .env
```

Skip copying if you already have `.env`. Set an API key from your dashboard:

```dotenv
DAYTONA_API_KEY=your_api_key
```

Optional `DAYTONA_API_URL` and `DAYTONA_TARGET` values map to
`Config.Endpoint` and `Config.Region`. When omitted, the SDK uses its normal
defaults. See [.env.example](.env.example) and
[provider configuration](../../../sdks/go/providers/daytona/README.md).

## Run

```sh
go run .
```

Expected output:

```text
Daytona sandbox created: <sandbox-id>
```

[main.go](main.go) loads `.env`, constructs the common client, and calls
`client.Create(ctx, nil)`. It propagates creation and SDK cleanup errors.
Shell variables take precedence over file values. The SDK itself does not load
`.env`.

The program provisions one real sandbox. Client cleanup does not stop/delete it;
use Daytona's tools for sandbox cleanup. Credentials stay in ignored `.env`.

## Verify locally

```sh
GOWORK=off go test ./...
```

This compiles the example without running creation.
For creation policies and SDK request tests, see
[create-daytona-with-policies](../create-daytona-with-policies/README.md).

The example leaves policies at provider defaults. Explicit zero semantics differ
by action; see [creation policies](../../../docs/sandbox-creation.md#policy-semantics).
Execution and lifecycle operations are outside this example.
