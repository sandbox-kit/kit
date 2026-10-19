# Create a Daytona sandbox with nested virtualization

Create from an image with `isolation.nestedVirtualization = true` (mapped to Daytona `kvm`).

## What this means

Nested virtualization is optional and account-dependent. Kit does not invent an isolation class; Daytona still chooses the sandbox class from the source.

## Edge cases

- Isolation `kind` / `providerKind` overrides are unsupported for Daytona.
- Org quotas and hardware availability remain Daytona-enforced.

## Run

```sh
cd examples/typescript/daytona/create-daytona-with-kvm
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Uses `await using`
so local client cleanup runs on success or failure.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
