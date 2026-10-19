# Create a Modal sandbox with idle termination

Create with **`idleTerminate` after 60 seconds** of inactivity.

## What this means

Modal idle timeouts automatically terminate a sandbox after inactivity (no
active exec/tunnels/etc., per Modal's idle rules). Separate from the maximum
lifetime `timeout`. See
[Sandboxes → Idle Timeouts](https://modal.com/docs/guide/sandboxes).

## Edge cases

- Delay must be a positive whole number of seconds when using `After`.
- The sandbox may disappear soon after create if it goes idle—expected.
- `Disabled` turns idle terminate off; zero delay is not “disabled.”

## Run

```sh
cd examples/typescript/modal/create-modal-idle
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
MODAL_TOKEN_ID=
MODAL_TOKEN_SECRET=
MODAL_APP_NAME=
MODAL_ENVIRONMENT=main
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/modal/create-modal-idle/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
