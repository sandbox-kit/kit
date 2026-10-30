# Create a Modal sandbox with exposed ports

Create with TLS (443) and plain (8080) port bindings.

## What this means

Each port needs an explicit `PortTransport`. Kit maps TLS → encrypted ports, plain → unencrypted ports, and HTTP/2 TLS → h2 ports.

## Edge cases

- Omitting `transport` fails locally.
- Public HTTP access is expressed through ports, not `security.publicAccess`.

## Run

```sh
cd examples/typescript/modal/create-modal-with-ports
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
MODAL_TOKEN_ID=
MODAL_TOKEN_SECRET=
MODAL_APP_NAME=
MODAL_ENVIRONMENT=main
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
