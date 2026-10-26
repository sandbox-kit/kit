# Sandbox Kit wire messages

This module is installed automatically with the sandbox client. The client
re-exports the messages a program constructs, including Config, CreateOptions,
and ErrorInfo.

Install the public client and one provider:

```sh
go get github.com/sandbox-kit/kit/sdks/go/sandbox@latest
go get github.com/sandbox-kit/kit/sdks/go/providers/daytona@latest
# or
go get github.com/sandbox-kit/kit/sdks/go/providers/modal@latest
```

Start with [the sandbox client](https://pkg.go.dev/github.com/sandbox-kit/kit/sdks/go/sandbox).

## License

MIT
