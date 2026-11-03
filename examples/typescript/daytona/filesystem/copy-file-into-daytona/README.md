# Copy a file into Daytona

Copy `DAYTONA_LOCAL_PATH` on this machine to `DAYTONA_REMOTE_PATH` in the sandbox. The program loads that sandbox with `client.get`.

## What this means

The call copies `DAYTONA_LOCAL_PATH` on this machine to `DAYTONA_REMOTE_PATH` in the sandbox. Missing parent directories of the remote file are created.

## Edge cases

- The sandbox path must be absolute. A blank local path is rejected before Daytona is called. The error field is `local_path`.
- A missing local file is not found before Daytona is called. The error field is `local_path`. A directory at that path is an invalid argument with the same field.
- An existing remote file is replaced.

## Run

```sh
cd examples/typescript/daytona/filesystem/copy-file-into-daytona
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
DAYTONA_LOCAL_PATH=./note.txt

```

## Warnings

- Overwrites the remote file when it already exists.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/daytona/filesystem/copy-file-into-daytona/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
