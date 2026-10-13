# Stop a Daytona sandbox

Retrieve one existing Daytona sandbox, check its resolved stop capability, and
stop it. Stop applies to Daytona containers, Linux VMs, Windows sandboxes, and
GPU sandboxes. Stopping is destructive for ephemeral sandboxes because Daytona
deletes them after they stop.

## Run

```sh
cd examples/go/daytona/stop-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. `Client.Close` does not perform another lifecycle operation.

```sh
GOWORK=off go test ./...
```
