module github.com/nexarise/microservices-factory/inventory

go 1.24

replace github.com/nexarise/microservices-factory/shared => ../shared

require (
	github.com/nexarise/microservices-factory/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.5
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
)
