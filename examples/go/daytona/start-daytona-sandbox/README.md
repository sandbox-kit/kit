# Start a Daytona sandbox

Retrieve one existing retained sandbox, check start support, and start it. Start
is exposed for non-terminal Daytona containers, Linux VMs, Windows sandboxes,
and GPU sandboxes. Daytona still enforces valid source states and account
permissions.

## Run

```sh
cd examples/go/daytona/start-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. `Client.Close` does not stop or delete the started sandbox.

```sh
GOWORK=off go test ./...
```
