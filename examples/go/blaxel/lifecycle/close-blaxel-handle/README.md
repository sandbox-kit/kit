# Close a Blaxel sandbox handle

Close one Blaxel sandbox handle without contacting Blaxel or changing the sandbox.

## What this means

`Close` releases local handle resources only. The sandbox keeps running until TTL, idle policy, or delete.

## Edge cases

- After close, further handle calls fail locally.
- Client close also does not delete sandboxes.

## Run

```sh
cd examples/go/blaxel/lifecycle/close-blaxel-handle
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
BL_SANDBOX_ID=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/lifecycle/close-blaxel-handle/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
