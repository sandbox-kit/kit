# Create a Daytona sandbox with TypeScript

This example creates a Daytona sandbox from the provider's default snapshot.
It prints the sandbox ID and leaves the cloud sandbox running. Closing the
client only releases local SDK resources. The example uses `await using`, so
local cleanup runs whether sandbox creation succeeds or throws.

```sh
cp .env.example .env
# Fill in DAYTONA_API_KEY.
npm install
npm run build
npm start
```

An uncaught SDK error is printed by Node.js and causes the process to exit with
a failure status; the example does not need to manage `process.exitCode`.

The Daytona SDK can also resolve its supported authentication environment
variables directly when `auth` is omitted from the shared client configuration.

## Authentication

The runnable example uses Daytona's standard environment authentication. Set
`DAYTONA_API_KEY` in `.env`. `DAYTONA_API_URL` and `DAYTONA_TARGET` are optional.
For an explicitly configured client, pass shared credentials instead:

```ts
const client = new Client({
  provider: daytona(),
  auth: {
    credentials: "apiKey",
    apiKey: { key: "your-api-key" },
  },
  region: "us",
});
```

Bearer authentication is also mapped through `auth.bearerToken`; it requires
`scope.organizationId`.

## Daytona representation notes

- Daytona selects the sandbox isolation class from the image or snapshot. The
  portable isolation kind must therefore be `IsolationKind.Default`;
  `nestedVirtualization` maps to Daytona's `kvm` setting.
- Creation-time resource overrides are supported only with an image source.
  Memory and disk use MiB in Sandbox Kit and must convert exactly to Daytona's
  whole-GiB values.
- Idle stop, idle pause, stopped archive, stopped delete, maximum lifetime, and
  ephemeral behavior are validated and converted from milliseconds to Daytona
  minutes.
- Outbound block/allow lists, linked sandboxes, outbound proxies, private/public
  access, environment variables, labels, and Daytona secret references are
  mapped when Daytona supports them.
- Unsupported portable intent is rejected before an API call instead of being
  silently ignored.
