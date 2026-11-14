# Create a Blaxel sandbox

Create one Blaxel sandbox from a registry image with an explicit name. Prints the sandbox ID and leaves the sandbox running.

## What this means

Create needs a Blaxel-valid `name` (lowercase alphanumeric with single hyphens,
at most 49 characters), an image, and client `region` (`BL_REGION`, default
`us-pdx-1`). If the name already exists, you get `AlreadyExists`. Closing the
client does not delete the sandbox.

## Edge cases

- Omitting `name` or `source.image`, or using an invalid name, fails before Blaxel is called.
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
