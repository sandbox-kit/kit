# Read a Daytona file

Read the bytes at `DAYTONA_REMOTE_PATH` with `Filesystem.ReadFile` / `sandbox.filesystem.readFile`. The program loads that sandbox with `client.get`.

## What this means

The call returns the file contents as bytes. Decode text in your own program. The sandbox path is absolute.

## Edge cases

- A relative path, a blank path, or a path that contains a NUL byte is rejected before Daytona is called. The error field is `path`.
- A missing file comes back from Daytona as not found.
- A directory path fails at the provider.

## Run

```sh
cd examples/typescript/daytona/read-daytona-file
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

- Needs an existing sandbox and a file at `DAYTONA_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/read-daytona-file/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
