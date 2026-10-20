# Create a Modal sandbox

Authenticate with a Modal **token pair**, scope to an **existing app**, and create
one sandbox from a registry **image**, with runtime, isolation, resources,
network, lifetime, and labels.

## What this means

Modal Sandboxes are containers you can create, run commands in, and terminate from
code. Create needs an **App** and an **Image**. Kit looks up an existing app; it
does not create apps. See Modal's
[Sandboxes guide](https://modal.com/docs/guide/sandboxes).

Closing the Kit client does not terminate the sandbox. Reclaim it with idle/max
timeouts or a terminate example.

## Edge cases

- Token display name, token id, and app name are different objects.
- Wrong/missing app fails at create time.
- Disk size cannot be raised on Sandboxes at create (see the error example).
- Default Modal resource request when omitted: 0.125 CPU cores and 128 MiB.

## Run

```sh
cd examples/typescript/modal/create-modal-sandbox
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

## Create an app if needed

Kit looks up an existing app; it does not create one. With the Modal Python CLI:

```sh
python3 -m pip install modal
modal token new
python3 -c 'import modal; modal.App.lookup("sandbox-kit-example", create_if_missing=True, environment_name="main")'
```

Set `MODAL_APP_NAME=sandbox-kit-example` in `.env`. See
[App.lookup](https://modal.com/docs/sdk/py/latest/App#lookup) and
[environments](https://modal.com/docs/guide/environments).

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/modal/create-modal-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
