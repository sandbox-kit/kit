PROTOC ?= protoc
GO ?= go
BIN := $(CURDIR)/work/bin

.PHONY: generate test

# Both generators are built from the version locked in tooling/go.mod.
# Stage 1 bootstraps the custom option types needed to build our plugin.
generate:
	mkdir -p $(BIN)
	cd tooling && GOWORK=off $(GO) build -o $(BIN)/protoc-gen-go google.golang.org/protobuf/cmd/protoc-gen-go
	$(PROTOC) -I proto --plugin=protoc-gen-go=$(BIN)/protoc-gen-go --go_out=tooling --go_opt=module=github.com/sandbox-kit/kit/tooling kit/codegen/v1/options.proto
	cd tooling && GOWORK=off $(GO) build -o $(BIN)/protoc-gen-kit-go ./cmd/protoc-gen-kit-go
	$(PROTOC) -I proto --plugin=protoc-gen-kit-go=$(BIN)/protoc-gen-kit-go --kit-go_out=core --kit-go_opt=module=github.com/sandbox-kit/kit/core kit/core/v1/client.proto
	$(PROTOC) -I proto --plugin=protoc-gen-kit-go=$(BIN)/protoc-gen-kit-go --kit-go_out=adapters/modal --kit-go_opt=module=github.com/sandbox-kit/kit/adapters/modal kit/adapters/modal/v1/adapter.proto
	$(PROTOC) -I proto --plugin=protoc-gen-kit-go=$(BIN)/protoc-gen-kit-go --kit-go_out=adapters/daytona --kit-go_opt=module=github.com/sandbox-kit/kit/adapters/daytona kit/adapters/daytona/v1/adapter.proto

test:
	cd core && GOWORK=off $(GO) test ./...
	cd tooling && GOWORK=off $(GO) test ./...
	cd adapters/modal && GOWORK=off $(GO) test ./...
	cd adapters/daytona && GOWORK=off $(GO) test ./...
	cd examples && GOWORK=off $(GO) test ./... && GOWORK=off $(GO) run . --help
