# Lino Redis

High level Redis wrapper for Go.

## Architecture

```plaintext
LinoRedis
    -> Fork(subPath) -> LinoRedis
    -> Get(subPath)
    -> Set(subPath, value)
    -> AnyCommand(subPath, ...)
    -> NewRedisItem(subPath) -> RedisItem
    -> NewAnyComplexType(subPath) -> AnyComplexType

RedisItem
    -> Get()
    -> Set(value)
    -> AnyCommand(...)
```

## Develop

start redis

`docker run --name lino_redis-test --rm -p 6379:6379 redis`
