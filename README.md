## command for generating the proto stubs etc

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ecommerce/ordermanagement/ordermanagement.proto

## Get dependencies

go mod tidy


