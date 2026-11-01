# Stat a Daytona path

Print the path, type, and size of `DAYTONA_REMOTE_PATH`. The program loads that sandbox with `client.get`.

## What this means

The call prints the path, type, and size. Daytona reports a file or a directory.

## Edge cases

- A relative path is rejected before Daytona is called.
- A missing path comes back from Daytona as not found.

## Run

```sh
cd examples/typescript/daytona/stat-daytona-file
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
DAYTONA_REMOTE_PATH=/tmp/note.txt

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox and a path at `DAYTONA_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/stat-daytona-file/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
