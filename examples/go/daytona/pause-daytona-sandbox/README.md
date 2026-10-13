# Pause a Daytona VM

Retrieve an existing Linux VM or Windows sandbox, check pause support, and
request a memory-preserving pause. Containers and GPU sandboxes are rejected
before the pause request. Choose a VM-capable snapshot, such as
`daytona-vm-small`, when creating the sandbox.

## Run

```sh
cd examples/go/daytona/pause-daytona-sandbox
cp .env.example .env
# Set DAYTONA_API_KEY and DAYTONA_SANDBOX_ID, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. `Client.Close` does not resume, stop, or delete the sandbox.

```sh
GOWORK=off go test ./...
```
