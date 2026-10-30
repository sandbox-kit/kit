# Create a Daytona sandbox with an injected secret

Create from an image and inject an existing Daytona secret as `GITHUB_TOKEN`.

## What this means

Daytona secrets require an environment variable name and typically egress-placeholder injection. The secret must already exist in your Daytona organization.

## Edge cases

- `DAYTONA_SECRET_NAME` must name a secret that already exists.
- Duplicate environment variables in `security.secrets` are rejected locally.

## Run

```sh
cd examples/typescript/daytona/create-daytona-with-secrets
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
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
DAYTONA_SECRET_NAME=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
