# Create a Modal sandbox

This standalone Go project authenticates with Modal through Sandbox Kit, creates
one sandbox from `alpine:3.21`, and prints its ID. It imports only the public SDK
and the Modal provider. Requires Go 1.26.1+.

## 1. Configure credentials and app

From the repository root:

```sh
cd examples/go/create-modal-sandbox
cp .env.example .env   # Skip if .env already exists.
```

Edit `.env`:

```dotenv
MODAL_TOKEN_ID=your_token_id
MODAL_TOKEN_SECRET=your_token_secret
MODAL_APP_NAME=your_existing_app_name
MODAL_ENVIRONMENT=main
```

Get the API token pair from **workspace settings → API tokens** in the Modal
dashboard. A token's display name is not its ID and is not the app name.
The dashboard's Secrets section is for values injected into containers;
API credentials are configured separately. See [Modal authentication](https://modal.com/docs/sdk/js/latest/intro)
and [Secrets](https://modal.com/docs/guide/secrets).

`MODAL_APP_NAME` is the name of an existing app in the environment selected by
`MODAL_ENVIRONMENT`. The example defaults to `main` if environment is unset/empty.
Environment names are visible in the dashboard's environment dropdown.
[Environment guide](https://modal.com/docs/guide/environments)

### Create an app if needed

Modal documents lookup with `create_if_missing=True` for creating an app that
owns sandboxes. Run this one-time setup with its Python SDK:

```sh
python3 -m pip install modal
modal token new
python3 -c 'import modal; modal.App.lookup("sandbox-kit-example", create_if_missing=True, environment_name="main")'
```

`modal token new` authenticates your local CLI. You can instead use credentials
already configured for that CLI. This command creates the app in the workspace
associated with those credentials; use the same workspace as the Go example's token.
[App lookup reference](https://modal.com/docs/sdk/py/latest/App#lookup)

Then set:

```dotenv
MODAL_APP_NAME=sandbox-kit-example
MODAL_ENVIRONMENT=main
```

Open the app's dashboard with `modal app dashboard sandbox-kit-example --env main`.
[CLI reference](https://modal.com/docs/cli/latest/app)

## 2. Run

```sh
go run .
```

Expected output:

```text
Modal sandbox created: <sandbox-id>
```

The program loads `.env`, passes credentials/app/environment through
`sandbox.Config`, and calls `client.Create` with an explicit registry image.
Existing shell environment variables take precedence over `.env`.
See [main.go](main.go) for the complete error-handled example and [.env.example](.env.example)
for the blank template. The SDK itself does not load `.env`.

## Scope and cleanup

The creation flow uses the image from [Modal's Go example](https://modal.com/docs/guide/sdk-javascript-go).
Sandbox Kit currently requires the app to already exist; it does not expose
Modal's create-app-if-missing option. This example does not execute commands,
mount volumes, or run sandbox lifecycle operations.

Running the program provisions a real sandbox. Deferred `client.Close` releases
SDK resources but does not terminate the sandbox. Manage sandbox cleanup through
Modal's dashboard/tools. `.env` is Git-ignored; `.env.example` has no credentials.

Compile without provisioning:

```sh
GOWORK=off go test ./...
```
