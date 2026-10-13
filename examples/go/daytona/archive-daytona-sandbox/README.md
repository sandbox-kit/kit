# Archive a Daytona container

Retrieve one existing Daytona sandbox, check its resolved capabilities, and
archive its retained filesystem. Archive is available for non-GPU containers;
Linux VMs, Windows sandboxes, and GPU sandboxes are rejected before the archive
request.

## Run

```sh
cd examples/go/daytona/archive-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. The sandbox class comes from the snapshot or image used at creation.
Closing the client does not stop, archive, or delete the sandbox.

```sh
GOWORK=off go test ./...
```
