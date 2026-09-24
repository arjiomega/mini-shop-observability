package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type ProductCache interface {
	Get(ctx context.Context, id int) (Product, error)
	Set(ctx context.Context, product Product, ttl time.Duration) error
	Delete(ctx context.Context, id int) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if err := redisotel.InstrumentTracing(client); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

func productCacheKey(id int) string {
	return fmt.Sprintf("product:%d", id)
}

func (c *RedisCache) Get(ctx context.Context, id int) (Product, error) {
	key := productCacheKey(id)

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return Product{}, err
	}

	var product Product

	if err := json.Unmarshal([]byte(value), &product); err != nil {
		return Product{}, err
	}

	return product, nil
}

func (c *RedisCache) Set(
	ctx context.Context,
	product Product,
	ttl time.Duration,
) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	key := productCacheKey(product.ID)

	return c.client.Set(
		ctx,
		key,
		data,
		ttl,
	).Err()
}

func (c *RedisCache) Delete(ctx context.Context, id int) error {
	return c.client.Del(
		ctx,
		productCacheKey(id),
	).Err()
}
