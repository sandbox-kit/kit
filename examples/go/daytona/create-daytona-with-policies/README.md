# Create a Daytona sandbox with policies

Create one sandbox from the default snapshot, disable auto-pause, and request
deletion ten minutes after stopping. This standalone module uses Sandbox Kit
and only the Daytona provider.

## Run with your account

Requires Go 1.26.9 or newer and this checkout. From the repository root:

```sh
cd examples/go/daytona/create-daytona-with-policies
cp .env.example .env
# Fill in DAYTONA_API_KEY.
go run .
```

Skip copying if `.env` already exists. Optional `DAYTONA_API_URL` and
`DAYTONA_TARGET` override endpoint/target. Shell variables take precedence
over `.env`, and credentials stay Git-ignored.

The program creates a real sandbox and prints its ID. It supplies policies at
creation; it does not call stop or delete. The deletion delay starts when the
sandbox stops. Closing the client only releases SDK resources.

## Verify without credentials

The example prints identity/resource origins and demonstrates structured common
errors with `errors.As`. Rejected policies expose stable kinds and field paths
before a creation request is sent. See [contracts](../../../../docs/responses-and-errors.md).

```sh
GOWORK=off go test -v ./...
```

Tests use dummy credentials and a local fake Daytona API. They run the actual
example through the public client and official SDK:

| Case                    | Expected behavior                                         |
| ----------------------- | --------------------------------------------------------- |
| Example creation        | Sends `autoPauseInterval: 0` and `autoDeleteInterval: 10` |
| Unspecified auto-stop   | Remains absent from the request                           |
| Disabled archive/delete | Rejected before a creation call                           |
| Immediate stop/pause    | Rejected before a creation call                           |
| Immediate delete        | Sends an explicit zero, rather than omitting it           |

Tests use a temporary configuration file, not your `.env`, and require
permission to bind a local port. They verify request mapping and rejection;
provider scheduling and real lifecycle behavior require live verification.

See [main.go](main.go), [tests](main_test.go), and
[creation policies](../../../../docs/sandbox-creation.md#policy-semantics).

## Error handling

`main.go` configures the provider and creates the sandbox. `errors.go` handles
application errors by extracting `*sandbox.Error` and switching on its `Kind`.
Ordinary `.env` errors are printed normally. Cleanup errors use the same handler
and do not replace a creation failure.

See the [provider error example](../handle-daytona-errors/README.md).
