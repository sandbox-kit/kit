# Create a Daytona sandbox with nested virtualization

Create from an image with `Isolation.NestedVirtualization` set to true.
Daytona maps that flag to KVM.

## What this means

Nested virtualization is optional and depends on the account. Daytona still
chooses the sandbox class from the source image.

## Edge cases

- Isolation kind overrides are rejected for Daytona.
- Organization quotas and hardware availability stay with Daytona.

## Run

```sh
cd examples/go/daytona/create-daytona-with-kvm
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Daytona.

## Related

- [typescript twin](../../../typescript/daytona/create-daytona-with-kvm/README.md)

- Catalog: [examples/go/README.md](../../README.md)
