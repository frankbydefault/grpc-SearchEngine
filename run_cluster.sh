docker run\
    -v "$PWD"/redis_cluster/7000/redis.conf:/usr/local/etc/redis/redis.conf\
    --rm\
    -d\
    --net=host\
    --name myredis-0\
    redis:6.2-alpine3.16\
    redis-server\
    /usr/local/etc/redis/redis.conf

docker run\
    -v "$PWD"/redis_cluster/7001/redis.conf:/usr/local/etc/redis/redis.conf\
    --rm\
    -d\
    --net=host\
    --name myredis-1\
    redis:6.2-alpine3.16\
    redis-server\
    /usr/local/etc/redis/redis.conf

docker run\
    -v "$PWD"/redis_cluster/7002/redis.conf:/usr/local/etc/redis/redis.conf\
    --rm\
    -d\
    --net=host\
    --name myredis-2\
    redis:6.2-alpine3.16\
    redis-server\
    /usr/local/etc/redis/redis.conf

docker run\
    -v "$PWD"/redis_cluster/7003/redis.conf:/usr/local/etc/redis/redis.conf\
    --rm\
    -d\
    --net=host\
    --name myredis-3\
    redis:6.2-alpine3.16\
    redis-server\
    /usr/local/etc/redis/redis.conf

docker run\
    -v "$PWD"/redis_cluster/7004/redis.conf:/usr/local/etc/redis/redis.conf\
    --rm\
    -d\
    --net=host\
    --name myredis-4\
    redis:6.2-alpine3.16\
    redis-server\
    /usr/local/etc/redis/redis.conf

docker run\
    -v "$PWD"/redis_cluster/7005/redis.conf:/usr/local/etc/redis/redis.conf\
    --rm\
    -d\
    --net=host\
    --name myredis-5\
    redis:6.2-alpine3.16\
    redis-server\
    /usr/local/etc/redis/redis.conf