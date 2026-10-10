# Create a linked Daytona sandbox

Create a parent sandbox from the default snapshot, then an ephemeral child
linked to that parent. Daytona schedules the child on the same runner.

## Run

Requires Go 1.26.1 or newer and this checkout.

```sh
cd examples/go/daytona/create-daytona-linked
cp .env.example .env
# Set DAYTONA_API_KEY, then:
go run .
```

The program creates two sandboxes. The child is ephemeral. Closing the client
does not stop either sandbox.

```sh
GOWORK=off go test ./...
```
