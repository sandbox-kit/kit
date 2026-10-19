# Create a Modal sandbox with outbound allowlists

Create with restricted egress, CIDR/domain allowlists, and private networking (`i6pn`).

## What this means

Complementary to the managed-proxy example: this sets allowlists and private network directly on the create request.

## Edge cases

- Empty inbound CIDR deny-all is not representable and is rejected.
- Linked sandbox / outbound proxy URL fields are unsupported on Modal.

## Run

```sh
cd examples/typescript/modal/create-modal-with-allowlists
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

- Catalog: [examples/typescript/README.md](../../README.md)
