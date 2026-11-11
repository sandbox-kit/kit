# Delete a Blaxel sandbox

Delete an existing Blaxel sandbox by name.

## What this means

`Delete` permanently removes the Blaxel sandbox. After success, `Info().State`
is `Deleted` and `Terminal` is true. This example uses `Wait: true` so removal
finishes before the call returns. Omit `Wait` to return when Blaxel accepts
the delete — recreating the same name too soon may return `Conflict`. Start,
stop, pause, archive, and terminate are not available.

## Edge cases

- Deletion destroys the sandbox and its filesystem.
- Closing the handle is not deletion.

## Run

```sh
cd examples/go/blaxel/lifecycle/delete-blaxel-sandbox
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
BL_API_KEY=
BL_WORKSPACE=
# Sandbox name (Blaxel id is the name).
BL_SANDBOX_NAME=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/lifecycle/delete-blaxel-sandbox/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
