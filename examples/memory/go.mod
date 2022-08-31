module github.com/mcandeia/dapr-components-go-sdk/examples/memory

go 1.19

require (
	github.com/dapr/components-contrib v1.8.2
	github.com/dapr/kit v0.0.2
	github.com/mcandeia/dapr-components-go-sdk v0.0.0-20220828173552-39ff56f2324d
)

require (
	github.com/dapr/dapr v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/net v0.0.0-20220630215102-69896b714898 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220622171453-ea41d75dfa0f // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace github.com/dapr/dapr => github.com/mcandeia/dapr v0.0.0-20220831143640-efe8777979fa

replace github.com/mcandeia/dapr-components-go-sdk => ../../.
