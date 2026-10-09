# Create a Daytona sandbox with policies

This standalone project creates one sandbox from Daytona's default snapshot,
explicitly disables auto-pause, and requests deletion ten minutes after stopping.
It uses the common Sandbox Kit API and installs only the Daytona provider.
Requires Go 1.26.1+.

## Run with your account

From the repository root:

```sh
cd examples/go/create-daytona-with-policies
cp .env.example .env
# Fill in DAYTONA_API_KEY from your dashboard.
go run .
```

Optional `DAYTONA_API_URL` and `DAYTONA_TARGET` override the endpoint/target.
Shell environment variables take precedence over `.env`.

Running provisions one real sandbox and prints its ID. The example requests
policies at creation; it does not itself stop or delete the sandbox. The delete
delay starts after the sandbox stops. Closing the client only releases SDK
resources. `.env` is ignored by Git.

## Verify without credentials

```sh
GOWORK=off go test -v ./...
```

Tests use dummy credentials and a local fake Daytona API. They run the example's
actual `run` function through the public client and official SDK, inspecting the
serialized creation request for `autoPauseInterval: 0` and
`autoDeleteInterval: 10`. They also assert that unspecified auto-stop stays absent.

Invalid disabled archive/delete and immediate stop/pause requests are exercised
through `client.Create` and must fail before a creation request is sent.
Immediate deletion is also tested for explicit zero serialization.
No real credentials or cloud sandboxes are used. Tests require permission to bind a local
port. Successful local transport verifies SDK request mapping, not provider-side
scheduling or lifecycle behavior.

See [main.go](main.go), [tests](main_test.go), and
[policy semantics](../../../docs/sandbox-creation.md).
