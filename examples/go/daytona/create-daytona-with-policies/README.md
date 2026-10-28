# Create a Daytona sandbox with lifetime policies

Create with **idle pause disabled** and **stopped-delete after 10 minutes**, then
print the sandbox id and whether resources were reported.

## What this means

Daytona can manage sandboxes with inactivity and stopped-time intervals:

| Policy | Effect (service guide) |
| ------ | ---------------------- |
| Auto-stop | Stop after N minutes idle; `0` disables (default often 15) |
| Auto-pause | Pause VM sandboxes after idle; mutually exclusive with auto-stop |
| Auto-archive | Archive stopped containers; `0` means max interval (30 days), not “now” |
| Auto-delete | Delete after stopped; `0` = immediate on stop; negative disables |
| TTL | Wall-clock destroy regardless of state |

This example maps Kit `idlePause: Disabled` and `stoppedDelete: After(10m)`.
Delays use whole minutes. See
[automated lifecycle management](https://www.daytona.io/docs/en/sandboxes/) and
Kit's [policy matrix](../../../../docs/sandbox-creation.md#policy-semantics)
(Kit follows the service guide where SDK comments disagree).

## Edge cases

- `PolicyModeAfter` with a **zero** delay means immediate, not disabled.
- Immediate stop/pause/archive and disabled archive/delete at create are rejected
  by the current Kit mapping.
- Positive stop and pause delays cannot both be set.
- Resources may be absent until the provider reports an allocation.

## Run

```sh
cd examples/go/daytona/create-daytona-with-policies
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/daytona/create-daytona-with-policies/README.md)

- Catalog: [examples/go/README.md](../../README.md)
