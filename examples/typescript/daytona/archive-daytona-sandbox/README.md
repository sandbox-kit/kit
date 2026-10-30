# Archive a Daytona sandbox

`Get` a sandbox, then archive it.

## What this means

Archive asks Daytona to move the sandbox filesystem to object storage. Daytona
accepts or rejects that sandbox. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- A rejected sandbox returns Daytona's error.
- Archive is not delete; recovery rules are Daytona-specific afterward.
- Mutates the target sandbox.

## Run

```sh
cd examples/typescript/daytona/archive-daytona-sandbox
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
- **Changes or destroys** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/archive-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
