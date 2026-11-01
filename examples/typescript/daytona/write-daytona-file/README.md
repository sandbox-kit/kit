# Write a Daytona file

Replace the file at `DAYTONA_REMOTE_PATH` with the bytes of `hello from sandbox-kit`. The program loads that sandbox with `client.get`.

## What this means

The call replaces the file with the bytes of `hello from sandbox-kit` and creates any missing parent directories.

## Edge cases

- A relative sandbox path is rejected before Daytona is called. The error field is `path`.
- An existing file at that path is replaced.
- A path that names a directory fails at the provider.

## Run

```sh
cd examples/typescript/daytona/write-daytona-file
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

- Overwrites the file at `DAYTONA_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/write-daytona-file/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
