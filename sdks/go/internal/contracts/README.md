# Sandbox Kit wire messages

Generated protobuf messages for modules under `sdks/go`. Only packages inside
`sdks/go` can import this module. Applications use the aliases on
[the sandbox client](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/sandbox)
(`Config`, `CreateOptions`, `ErrorInfo`, and related types).

Install the public client and one provider:

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona@latest
# or
go get github.com/sandbox-kit/kit/sdks/go/providers/modal@latest
```

## License

MIT
