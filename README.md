```sql
create database search_engine;  
create user searchengine@localhost identified by 'S34rch3r_3ng1n3';  
grant all privileges on search_engine.* to searchengine@localhost;  
```

Generate Protoc 
```
export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     protos/searchEng.proto

```

Run Server
`go run server/main.go 2>&1`

Run Client
`go run client/main.go --message=[search string]`