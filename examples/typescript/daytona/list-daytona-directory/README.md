# List a Daytona directory

Print the names in the directory at `DAYTONA_REMOTE_PATH`. The program loads that sandbox with `client.get`.

## What this means

The call prints the names in one directory. It does not walk into subdirectories.

## Edge cases

- A relative path is rejected before Daytona is called.
- A missing directory comes back from Daytona as not found.
- A file path fails at the provider.

## Run

```sh
cd examples/typescript/daytona/list-daytona-directory
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
DAYTONA_REMOTE_PATH=/tmp

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox and a directory at `DAYTONA_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/list-daytona-directory/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
