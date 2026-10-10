# Create an ephemeral Daytona sandbox

Create one sandbox from the default snapshot and mark it ephemeral. Daytona
deletes an ephemeral sandbox when it stops.

## Run

Requires Go 1.26.9 or newer and this checkout.

```sh
cd examples/go/daytona/create-daytona-ephemeral
cp .env.example .env
# Set DAYTONA_API_KEY, then:
go run .
```

Closing the client does not stop the sandbox. Stopping it deletes it.

```sh
GOWORK=off go test ./...
```
