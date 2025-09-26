module github.com/nimbodex/microservices-factory/inventory

go 1.24

replace github.com/nimbodex/microservices-factory/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/nimbodex/microservices-factory/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
	google.golang.org/grpc v1.69.2
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
