module github.com/sandbox-kit/kit/examples/go/blaxel/write-blaxel-file

go 1.26.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/sandbox-kit/kit/sdks/go/providers/blaxel v0.1.5
	github.com/sandbox-kit/kit/sdks/go/sandbox v0.1.5
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect
	buf.build/go/protovalidate v1.4.0 // indirect
	cel.dev/cel-go v0.32.0 // indirect
	cel.dev/expr v0.25.3 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/blaxel-ai/sdk-go v0.31.0 // indirect
	github.com/sandbox-kit/kit/sdks/go/internal/contracts v0.1.5 // indirect
	github.com/sandbox-kit/kit/sdks/go/internal/sdk v0.1.5 // indirect
	github.com/tidwall/gjson v1.19.0 // indirect
	github.com/tidwall/match v1.2.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/exp v0.0.0-20260820142414-ca536658362e // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260819154853-08b0e4226688 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260819154853-08b0e4226688 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sandbox-kit/kit/sdks/go/sandbox => ../../../../../sdks/go/sandbox

replace github.com/sandbox-kit/kit/sdks/go/providers/blaxel => ../../../../../sdks/go/providers/blaxel

replace github.com/sandbox-kit/kit/sdks/go/internal/contracts => ../../../../../sdks/go/internal/contracts

replace github.com/sandbox-kit/kit/sdks/go/internal/sdk => ../../../../../sdks/go/internal/sdk
