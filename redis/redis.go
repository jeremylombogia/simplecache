package redis

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v8"
	cache "github.com/sahalazain/simplecache"
)

const defaultNS = "redis"
const schema = "redis"

//Cache redis cache object
type Cache struct {
	client *redis.Client
	ns     string
}

func init() {
	cache.Register(schema, NewCache)
}

//NewCache create new redis cache
func NewCache(url *url.URL) (cache.Cache, error) {
	p, _ := url.User.Password()
	rClient := redis.NewClient(&redis.Options{
		Addr:     url.Host,
		Password: p,
		DB:       0, // use default DB
	})
	cache := &Cache{
		client: rClient,
		ns:     strings.TrimPrefix(url.Path, "/"),
	}
	_, err := cache.client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return cache, nil
}

//NewRedisCache creating instance of redis cache
func NewRedisCache(address string, ns string) (*Cache, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	cache := &Cache{
		client: rClient,
		ns:     ns,
	}
	_, err := cache.client.Ping(context.Background()).Result()
	//log.Debug(pong, err)
	return cache, err
}

//Set set value
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	switch value.(type) {
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, []byte:
		return c.client.Set(ctx, c.ns+key, value, time.Duration(expiration)*time.Second).Err()
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return c.client.Set(ctx, c.ns+key, b, time.Duration(expiration)*time.Second).Err()
	}
}

//Get get value
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(ctx, c.ns+key).Bytes()
}

//GetObject get object value
func (c *Cache) GetObject(ctx context.Context, key string, doc interface{}) error {
	b, err := c.client.Get(ctx, c.ns+key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, doc)
}

//GetString get string value
func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, c.ns+key).Result()
}

//GetInt get int value
func (c *Cache) GetInt(ctx context.Context, key string) (int64, error) {
	return c.client.Get(ctx, c.ns+key).Int64()
}

//GetFloat get float value
func (c *Cache) GetFloat(ctx context.Context, key string) (float64, error) {
	return c.client.Get(ctx, c.ns+key).Float64()
}

//Exist check if key exist
func (c *Cache) Exist(ctx context.Context, key string) bool {
	return c.client.Exists(ctx, c.ns+key).Val() > 0
}

//Delete delete record
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.ns+key).Err()
}

//RemainingTime get remaining time
func (c *Cache) RemainingTime(ctx context.Context, key string) int {
	return int(c.client.TTL(ctx, c.ns+key).Val().Seconds())
}

//Close close connection
func (c *Cache) Close() error {
	return c.client.Close()
}
