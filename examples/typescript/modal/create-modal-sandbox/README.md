# Create a Modal sandbox with TypeScript

This example uses the shared Sandbox Kit API to create one Modal gVisor sandbox
with explicit runtime, resource, network, and lifetime settings. Closing the
client releases local SDK resources; the sandbox terminates through its maximum
lifetime or an explicit `sandbox.terminate()` call.

```sh
cp .env.example .env
npm install
npm run build
npm start
```
