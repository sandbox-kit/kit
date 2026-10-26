# Create a Daytona sandbox with language settings

Create from `daytona-small` with toolbox **`language=python`** and
**`user=daytona`**, then print the resolved isolation kind.

## What this means

Daytona's `language` selects the programming-language runtime used for direct
code execution in the sandbox (`python`, `typescript`, or `javascript`). If
omitted, Daytona defaults to Python. `user` sets the OS user. See
[sandboxes → languages](https://www.daytona.io/docs/en/sandboxes/).

These settings configure the toolbox on top of the snapshot; they do not replace
the snapshot's sandbox class.

## Edge cases

- Unsupported language strings fail local validation.
- Printed isolation kind is provider-reported metadata after create.
- Entrypoint/workdir-style process control is not what Daytona `language` means.

## Run

```sh
cd examples/go/daytona/create-daytona-with-runtime
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

- [typescript twin](../../../typescript/daytona/create-daytona-with-runtime/README.md)

- Catalog: [examples/go/README.md](../../README.md)
