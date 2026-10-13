# Delete a Daytona sandbox

Retrieve one existing Daytona sandbox, check delete support, and wait for
irreversible deletion. Delete applies to Daytona containers, Linux VMs, Windows
sandboxes, and GPU sandboxes; provider permissions and current state still
determine whether Daytona accepts the request.

## Run

```sh
cd examples/go/daytona/delete-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. After success, the local handle is terminal and exposes no further
lifecycle capabilities. `Client.Close` only releases SDK resources.

```sh
GOWORK=off go test ./...
```
