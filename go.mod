module github.com/mcandeia/dapr-components-go-sdk

go 1.19

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220630215102-69896b714898 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220622171453-ea41d75dfa0f // indirect
)

require (
	github.com/dapr/components-contrib v1.8.2
	github.com/dapr/dapr v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

replace github.com/dapr/dapr => github.com/mcandeia/dapr v0.0.0-20220827174239-0397ce99b343
