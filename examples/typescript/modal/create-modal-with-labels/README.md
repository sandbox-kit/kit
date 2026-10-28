# Create a Modal sandbox with labels and environment

Create with process **environment** variables and **labels/tags**.

## What this means

Environment variables are injected into the sandbox process. Labels/tags help
organize sandboxes in the Modal UI and APIs. Kit validates environment key syntax
when mapping.

## Edge cases

- Do not put secrets in labels or example env values committed to git.
- Name and labels on the returned metadata are the values from the create request.

## Run

```sh
cd examples/typescript/modal/create-modal-with-labels
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

- [go twin](../../../go/modal/create-modal-with-labels/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
