# create-modal-idle with TypeScript

This independent example mirrors the Go `create-modal-idle` example through the Sandbox Kit TypeScript API.

```sh
cp .env.example .env
# Fill in the provider credentials and any required sandbox ID.
npm install
npm run build
npm start
```

Closing the client or sandbox handle only releases local resources unless this example explicitly invokes a lifecycle operation.

