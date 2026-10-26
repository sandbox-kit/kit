module github.com/sandbox-kit/kit/sdks/go/internal/sdk

go 1.26.0

require (
	github.com/sandbox-kit/kit/sdks/go/internal/contracts v0.1.1
	google.golang.org/protobuf v1.36.12
)

require buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect

replace github.com/sandbox-kit/kit/sdks/go/internal/contracts => ../contracts
