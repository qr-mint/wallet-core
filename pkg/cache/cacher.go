package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cacher interface {
	SetWithTTL(key string, value interface{}, expiration time.Duration) error
	Get(key string, value interface{}) error
}

type Params struct {
	Host     string
	Port     int
	Password string
	Database int
}

type redisCacher struct {
	client *redis.Client
}

func NewRedisCacher(params Params) Cacher {
	return &redisCacher{
		client: redis.NewClient(
			&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", params.Host, params.Port),
				Password: params.Password,
				DB:       params.Database,
			},
		),
	}
}

func (c redisCacher) SetWithTTL(key string, value interface{}, expiration time.Duration) error {
	_, err := c.client.Ping(context.Background()).Result()

	if err != nil {
		return fmt.Errorf("unable to connect to Redis: %v", err)
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error while trying marshal cache data: %s", err)
	}

	c.client.Set(context.Background(), key, data, expiration)

	return nil
}

func (c redisCacher) Get(key string, value interface{}) error {
	result, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	if result == "" {
		return nil
	}

	err = json.Unmarshal([]byte(result), value)
	if err != nil {
		return fmt.Errorf("error while trying unmarshal cache data: %s", err)
	}

	return nil
}
