# Copy a file from Daytona

Copy `DAYTONA_REMOTE_PATH` from the sandbox to `DAYTONA_LOCAL_PATH` on this machine. The program loads that sandbox with `client.get`.

## What this means

The call copies `DAYTONA_REMOTE_PATH` from the sandbox to `DAYTONA_LOCAL_PATH` on this machine.

## Edge cases

- The sandbox path must be absolute. A blank local path is rejected before Daytona is called. The error field is `local_path`.
- A missing remote file comes back from Daytona as not found.
- A missing local parent directory is created. A local directory is an invalid argument, and the field is `local_path`.
- An existing local file is replaced.

## Run

```sh
cd examples/typescript/daytona/copy-file-from-daytona
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

- Overwrites the local file when it already exists.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/copy-file-from-daytona/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
