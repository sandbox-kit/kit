# Create a Daytona sandbox from an image

Create from a public registry image (`ubuntu:22.04`) with explicit CPU, memory,
and disk.

## What this means

Image create builds/runs from a registry reference and lets you set resources.
Daytona documents resources in **vCPU and GiB**. Kit accepts MiB and converts
exactly to GiB (divide by 1024; remainders are rejected). CPU must be a whole
number of cores.

Documented organization defaults for image create are often 1 vCPU / 1 GiB /
3 GiB disk, with higher account maxima. See
[sandboxes → resources](https://www.daytona.io/docs/en/sandboxes/) and
[limits](https://www.daytona.io/docs/limits).

## Edge cases

- Valid Kit memory/disk examples: `1024`, `2048`, `4096` MiB. `1536` MiB fails
  local conversion.
- Image pull permissions and org quotas remain Daytona-enforced.
- Image references: lowercase paths, optional registry host/port, tag ≤128 chars
  or SHA-256 digest, max 512 chars.

## Run

```sh
cd examples/go/daytona/create-daytona-from-image
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
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

- [typescript twin](../../../typescript/daytona/create-daytona-from-image/README.md)

- Catalog: [examples/go/README.md](../../README.md)
