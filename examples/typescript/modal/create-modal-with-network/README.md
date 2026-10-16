# create-modal-with-network with TypeScript

This independent example creates a sandbox whose outbound traffic uses the
existing Modal Proxy named by `MODAL_PROXY_NAME`.

```sh
cp .env.example .env
# Fill in the provider credentials, app name, and existing proxy name.
npm install
npm run build
npm start
```

Closing the client or sandbox handle only releases local resources unless this example explicitly invokes a lifecycle operation.

