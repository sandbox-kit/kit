# Create a Daytona sandbox from an image

Create one sandbox from `ubuntu:22.04` with 1 vCPU, 1 GiB memory, and 3 GiB disk.

Resource overrides require an image. CPU must be a whole number of cores.
Memory and disk must be whole GiB, expressed as MiB.

## Run

Requires Go 1.26.1 or newer and this checkout.

```sh
cd examples/go/daytona/create-daytona-from-image
cp .env.example .env
# Set DAYTONA_API_KEY, then:
go run .
```

The image needs a tag or digest. Daytona rejects `latest`, `lts`, and `stable`.
Closing the client does not delete the sandbox.

```sh
GOWORK=off go test ./...
```
