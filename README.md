# simplecache
Simple cache abstraction for Golang

Supported cache
  - In Memory
  - Redis
  - LRU

## Quick Start

Installation
    $ go get github.com/sahalazain/simplecache


## Usage

Interface:
```
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetObject(ctx context.Context, key string, doc interface{}) error
	GetString(ctx context.Context, key string) (string, error)
	GetInt(ctx context.Context, key string) (int64, error)
	GetFloat(ctx context.Context, key string) (float64, error)
	Exist(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
	RemainingTime(ctx context.Context, key string) int
	Close() error
}
```

Example:

```go
package main

import (
	cache "github.com/sahalazain/simplecache"
	_ "github.com/sahalazain/simplecache/mem"
	_ "github.com/sahalazain/simplecache/redis"
	_ "github.com/sahalazain/simplecache/lru"
)

func main() {
	// Use in-memory store
	memcache, _ := cache.New("mem://")
	rediscache, _ := cache.New("redis://<user>:<pass>@localhost:6379/prefix")
	lru, _ := cache.New("lru://local/1024")
}
```
