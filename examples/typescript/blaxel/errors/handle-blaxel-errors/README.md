# Handle Blaxel validation errors

Demonstrate typed validation errors from the Blaxel provider without creating a sandbox.

## What this means

Requests an unsupported CPU override so Kit fails before calling Blaxel. Applications inspect `SandboxError` kind and field.

## Edge cases

- Does not provision cloud resources when validation fails as expected.
- Authentication still requires a configured client.

## Run

```sh
cd examples/typescript/blaxel/errors/handle-blaxel-errors
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Skip copying if
`.env` already exists.

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

- [go twin](../../../../go/blaxel/errors/handle-blaxel-errors/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
