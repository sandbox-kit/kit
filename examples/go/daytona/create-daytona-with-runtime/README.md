# Create a Daytona sandbox with runtime settings

Create one Daytona container from `daytona-small` with Python selected as the
Daytona toolbox language and `daytona` as the sandbox user.

Daytona maps `Runtime.Language` and `Runtime.User` during creation. Supported
toolbox languages are `python`, `javascript`, and `typescript`. Runtime version,
working directory, entrypoint, and PTY are not currently mapped and fail before
the provider call instead of being ignored.

The snapshot determines the isolation class. `daytona-small` creates a
container; a VM snapshot such as `daytona-vm-small` creates a Linux VM. Setting
`Isolation.Kind` does not override the snapshot class.

## Run

Requires Go 1.26.9 or newer and this checkout.

```sh
cd examples/go/daytona/create-daytona-with-runtime
cp .env.example .env
# Set DAYTONA_API_KEY, then:
go run .
```

The example constructs its own client with `sandbox.NewClient`. Optional
`DAYTONA_API_URL` and `DAYTONA_TARGET` values select the Daytona endpoint and
target. Closing the client does not stop or delete the sandbox.

```sh
GOWORK=off go test ./...
```
