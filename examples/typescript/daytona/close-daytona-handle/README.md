# Close a Daytona sandbox handle

`Get` a sandbox, then call **`Sandbox.Close` / `sandbox.close()`**.

## What this means

Closing the handle does not contact Daytona and does not stop, archive, or
delete the sandbox. Closing the client closes your connection to Daytona.
You can keep using the handle. The next call returns Daytona's result.

## Edge cases

- Safe against a sandbox you want to keep running or stopped.
- A later call on the handle goes to Daytona.
- Client close releases the SDK connection. A later call returns Daytona's result.

## Run

```sh
cd examples/typescript/daytona/close-daytona-handle
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
DAYTONA_SANDBOX_ID=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- Handle close does not change the cloud sandbox. Client close releases the SDK connection.

## Related

- [go twin](../../../go/daytona/close-daytona-handle/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
