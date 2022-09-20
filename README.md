# GoGRPC Searcher

Este proyecto utiliza como motor de base de datos [MySQL](https://www.mysql.com/).  
Para el almacenamiento de cache se utiliza [Redis](https://redis.io/), en un contenedor [Docker](https://www.docker.com/) con la imagen [redis:6.2-alpine3.16](https://hub.docker.com/_/redis).  
Para la comunicación entre los servicios se utiliza [gRPC](https://grpc.io/), con su implementación en [Go](https://go.dev/).  
Para la lectura de la información y su posterior almacenamiento en la base de datos se utiliza [Python](https://www.python.org/), con la librería [pyDAL](https://github.com/web2py/pydal).  
Para la comunicación con Redis desde Go se utilizó [go-redis](https://github.com/go-redis/redis) y para la comunicación con MySQL desde Go se utilizó [go-sql-driver](https://github.com/go-sql-driver/mysql).

## Preparando el servicio

* Para crear la base de datos, esto se hará con el motor MySQL.

```sql
create database search_engine;  
create user searchengine@localhost identified by 'S34rch3r_3ng1n3';  
grant all privileges on search_engine.* to searchengine@localhost;  
```

Es posible utilizar el dump de la base de datos *searches.db* o utilizar el script *crawler.py* en *databases/* para llenar la base de datos.

* Compilar protoc
```
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     protos/searchEng.proto
```

* Iniciar contenedores del cluster de Redis
```bash
chmod +x ./run_cluster.sh
sudo ./run_cluster.sh
```

* Crear el cluster en los contenedores
```sh
sudo docker exec -it myredis-0 sh
redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 --cluster-replicas 1
```
*Si se solicita confirmación, escribir yes*

* Para detener las instancias de Redis
```bash
chmod +x ./stop_cluster.sh
sudo ./stop_cluster.sh
```

## Iniciando cliente y servidor

* Iniciar servidor
```bash
go run server/main.go
```

* Iniciar cliente
```bash
go run client/main.go
```

## Usando el cliente

Una vez se ejecuta el cliente, este solicitará un término de búsqueda.  
El cliente buscará en el cache de Redis, si no se encuentra, se buscará en la base de datos y se almacenará en el cache.

## Analizando cache
Para analizar el cache basta con utilizar el comando
```bash
redis-cli -c -p 7000
```
Y luego utilizar *get* seguido de la llave para ver dónde se encuentra el valor en el cluster y obtener el valor o *set* para agregar un valor al cache, indicando en qué posición del cluster se encuentra.
