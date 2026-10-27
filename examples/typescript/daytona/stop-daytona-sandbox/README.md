# Stop a Daytona sandbox

`Get` a sandbox, check **Stop**, then stop it.

## What this means

Stop asks Daytona to stop the sandbox. A retained sandbox can usually be started
again. An ephemeral sandbox can be deleted by its policy. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Daytona applies its own stop rules, including ephemeral deletion.
- Ephemeral + stop deletes the sandbox.
- **Mutates** the ID in `.env`.

## Run

```sh
cd examples/typescript/daytona/stop-daytona-sandbox
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/daytona/stop-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
