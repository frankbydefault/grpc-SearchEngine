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
```bash
go run server/main.go
```

Run Client
```bash
go run client/main.go --message=[search string]
```

To run cluster
```bash
chmod +x ./run_cluster.sh
sudo ./run_cluster.sh
```

To enable cluster
```sh
sudo docker exec -it myredis-0 sh
redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 --cluster-replicas 1
```
*When asked, type yes*

To stop cluster
```bash
chmod +x ./stop_cluster.sh
sudo ./stop_cluster.sh
```