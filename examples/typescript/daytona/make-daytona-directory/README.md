# Create a Daytona directory

Create the directory at `DAYTONA_REMOTE_PATH`, including any missing parents. The program loads that sandbox with `client.get`.

## What this means

The call creates `DAYTONA_REMOTE_PATH` and any missing parents.

## Edge cases

- A relative path is rejected before Daytona is called.
- Daytona returns an error when that directory already exists.

## Run

```sh
cd examples/typescript/daytona/make-daytona-directory
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
DAYTONA_REMOTE_PATH=/tmp/kit-files

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Creates the directory inside the existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/make-daytona-directory/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
