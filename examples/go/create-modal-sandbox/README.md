# Create a Modal sandbox

Authenticate with Modal, create one sandbox from `alpine:3.21`, and print its ID.
This standalone Go module imports Sandbox Kit and only the Modal provider.

## Configure credentials and scope

Requires Go 1.26.1 or newer and this checkout. From the repository root:

```sh
cd examples/go/create-modal-sandbox
cp .env.example .env
```

Skip copying if you already have `.env`. Fill in:

```dotenv
MODAL_TOKEN_ID=your_token_id
MODAL_TOKEN_SECRET=your_token_secret
MODAL_APP_NAME=your_existing_app_name
MODAL_ENVIRONMENT=main
```

Create an API token pair in Modal workspace settings. The token's display name,
its ID, and the app name identify different things. Dashboard Secrets hold values
injected into workloads; they are separate from API credentials.
See [Go client configuration](https://modal.com/docs/sdk/go/latest/Client).

The app must exist in the same workspace/environment as the token.
The example defaults `MODAL_ENVIRONMENT` to `main` when unset or empty.
See [Modal environments](https://modal.com/docs/guide/environments).

## Create an app if needed

One option is the Modal Python SDK/CLI. Authenticate it to the same workspace,
then create the app once:

```sh
python3 -m pip install modal
modal token new
python3 -c 'import modal; modal.App.lookup("sandbox-kit-example", create_if_missing=True, environment_name="main")'
```

Use `MODAL_APP_NAME=sandbox-kit-example` in `.env`. Replace the name/environment
in that command if you want another app. Existing CLI credentials can be used
instead of `modal token new`.

See [App.lookup](https://modal.com/docs/sdk/py/latest/App#lookup).
Sandbox Kit's creation mapping looks up an existing app; it does not create one.

## Run

```sh
go run .
```

Expected output:

```text
Modal sandbox created: <sandbox-id>
```

[main.go](main.go) loads `.env`, constructs `sandbox.Config`, creates the image-based
sandbox, and propagates creation/cleanup errors. Shell variables take precedence
over file values. The SDK itself does not load `.env`.

The program creates a real sandbox. `client.Close` releases SDK resources without
terminating it; manage cleanup through Modal tools or the provider lifetime.
Credentials stay in ignored `.env`; [.env.example](.env.example) contains no secrets.

## Verify locally

```sh
GOWORK=off go test ./...
```

This compiles the example without running creation.
For resources and local rejection tests, see
[create-modal-with-resources](../create-modal-with-resources/README.md).
Execution, lifecycle operations, and storage are outside this example.

## Error handling

`main.go` configures the provider and creates the sandbox. `errors.go` handles
application errors by extracting `*sandbox.Error` and switching on its `Kind`.
Ordinary `.env` errors are printed normally. Cleanup errors use the same handler
and do not replace a creation failure.

See the [provider error example](../handle-modal-errors/README.md).
