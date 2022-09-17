Start server
go run server/server.go

Test client
go run client/client.go --name=Alice

gen PROTOC 

export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     protos/searchEng.proto