module github.com/Lempi-sudo/lempi-rocket-project/inventory

go 1.24.3

replace github.com/Lempi-sudo/lempi-rocket-project/shared => ../shared

require (
	github.com/Lempi-sudo/lempi-rocket-project/shared v0.0.0-00010101000000-000000000000
	github.com/samber/lo v1.51.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
)
