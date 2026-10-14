# Sandbox Kit for TypeScript

The public client owns the selected provider SDK. Close it explicitly when the
application is finished. Closing the client releases local SDK resources; it
does not stop or delete cloud sandboxes.

```ts
const client = new Client(config);

try {
  const sandbox = await client.create();
  console.log(sandbox.id);
} finally {
  await client.close();
}
```

Runtimes that support explicit resource management may use `await using` as an
optional convenience:

```ts
await using client = new Client(config);
const sandbox = await client.create();
```

The second form invokes `client[Symbol.asyncDispose]()` automatically at the end
of the scope. Both forms have identical cleanup semantics.
