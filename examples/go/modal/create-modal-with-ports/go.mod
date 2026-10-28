module github.com/sandbox-kit/kit/examples/go/modal/create-modal-with-ports

go 1.26.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/sandbox-kit/kit/sdks/go/providers/modal v0.1.2
	github.com/sandbox-kit/kit/sdks/go/sandbox v0.1.2
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect
	buf.build/go/protovalidate v1.4.0 // indirect
	cel.dev/cel-go v0.32.0 // indirect
	cel.dev/expr v0.25.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/aristanetworks/gomap v0.0.0-20230726210543-f4e41046dced // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kisielk/og-rek v1.3.0 // indirect
	github.com/modal-labs/modal-client/go v0.11.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sandbox-kit/kit/sdks/go/internal/contracts v0.1.2 // indirect
	github.com/sandbox-kit/kit/sdks/go/internal/sdk v0.1.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/exp v0.0.0-20260820142414-ca536658362e // indirect
	golang.org/x/net v0.60.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260819154853-08b0e4226688 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260819154853-08b0e4226688 // indirect
	google.golang.org/grpc v1.83.2 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/sandbox-kit/kit/sdks/go/sandbox => ../../../../sdks/go/sandbox

replace github.com/sandbox-kit/kit/sdks/go/providers/modal => ../../../../sdks/go/providers/modal

replace github.com/sandbox-kit/kit/sdks/go/internal/contracts => ../../../../sdks/go/internal/contracts

replace github.com/sandbox-kit/kit/sdks/go/internal/sdk => ../../../../sdks/go/internal/sdk
