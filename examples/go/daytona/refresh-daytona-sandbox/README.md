# Refresh Daytona metadata

Retrieve a sandbox handle and refresh its cached metadata from Daytona. Copy
the provider's latest state, isolation class, resources, and native metadata
into `instance.Info()`. Refresh is read-only and works for every Daytona
sandbox class.

## Run

```sh
cd examples/go/daytona/refresh-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. Neither refresh nor `Client.Close` changes the cloud sandbox.

```sh
GOWORK=off go test ./...
```
