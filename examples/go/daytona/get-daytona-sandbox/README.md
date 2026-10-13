# Get a Daytona sandbox

Retrieve current metadata for one existing sandbox and print its portable state,
native Daytona state, and class-aware lifecycle capabilities. The returned
metadata distinguishes containers, Linux VMs, Windows sandboxes, and GPU-backed
containers when Daytona supplies that information.

## Run

```sh
cd examples/go/daytona/get-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. `Get` and `Client.Close` do not change the cloud sandbox.

```sh
GOWORK=off go test ./...
```
