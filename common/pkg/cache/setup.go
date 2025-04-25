package cache

import (
	"context"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewCache(dsn string) *Redis {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &Redis{redis.NewClient(opt)}
}

func (c *Redis) FindOne(ctx context.Context, key string, target interface{}) error {
	jsonStr, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonStr), target)
}

func (c *Redis) InsertOne(ctx context.Context, key string, duration time.Duration, input interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, string(bytes), duration).Err()
}
