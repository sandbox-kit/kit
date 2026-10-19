# Create a Daytona sandbox with a maximum lifetime

Create from an image with `lifetime.maximumLifetime` set to 60 minutes.

## What this means

Daytona maps maximum lifetime to `ttlMinutes`. Kit requires a positive whole-minute duration.

## Edge cases

- Durations that are not whole minutes are rejected locally.
- Idle stop/pause and stopped archive/delete policies are covered by `create-daytona-with-policies`.

## Run

```sh
cd examples/typescript/daytona/create-daytona-with-ttl
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
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
