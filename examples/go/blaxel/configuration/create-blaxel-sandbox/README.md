# Create a Blaxel sandbox

Create one Blaxel sandbox from a registry image with an explicit name. Prints the sandbox ID and leaves the sandbox running.

## What this means

Blaxel create requires a nonblank `name` and an image source. Kit calls native create (not create-if-not-exists). Region comes from the client config (`BL_REGION`, default `us-pdx-1`). Closing the Kit client only drops the local SDK connection.

## Edge cases

- Omitting `name` or `source.image` fails before Blaxel is called.
- A create error is not proof that no sandbox was provisioned.
- Reclaim sandboxes with TTL / idle policies or `delete`.

## Run

```sh
cd examples/go/blaxel/configuration/create-blaxel-sandbox
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
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/configuration/create-blaxel-sandbox/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
