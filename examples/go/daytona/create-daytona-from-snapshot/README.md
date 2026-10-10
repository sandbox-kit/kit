# Create a Daytona sandbox from a snapshot

Create one sandbox from the stock `daytona-small` snapshot and print its ID.

## Run

Requires Go 1.26.1 or newer and this checkout.

```sh
cd examples/go/daytona/create-daytona-from-snapshot
cp .env.example .env
# Set DAYTONA_API_KEY, then:
go run .
```

`daytona-small` is a 1 vCPU container snapshot. Other stock names, such as
`daytona-medium` and `daytona-vm-small`, use the same `Snapshot` field.
Closing the client does not delete the sandbox.

```sh
GOWORK=off go test ./...
```
