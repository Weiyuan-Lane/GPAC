package rediswrapper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisWrapper struct {
	redisClient *redis.Client
}

func New(client *redis.Client) *RedisWrapper {
	return &RedisWrapper{
		redisClient: client,
	}
}

func (r *RedisWrapper) Get(key string) (*string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &val, nil
}

func (r *RedisWrapper) MultipleGet(keys []string) (map[string]string, error) {
	val, err := r.redisClient.MGet(ctx, keys...).Result()
	if err == redis.Nil {
		return map[string]string{}, nil
	} else if err != nil {
		return map[string]string{}, err
	}

	result := map[string]string{}
	for i, key := range keys {
		strPtr, ok := val[i].(*string)
		if ok && strPtr != nil {
			result[key] = *strPtr
		}
	}

	return result, nil
}

func (r *RedisWrapper) Set(key, val string, ttl int) error {
	secondDuration := time.Duration(ttl) * time.Second
	err := r.redisClient.Set(ctx, key, val, secondDuration).Err()
	return err
}

func (r *RedisWrapper) MultipleSet(valMap map[string]string, ttl int) error {
	for k, v := range valMap {
		err := r.Set(k, v, ttl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisWrapper) Delete(key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
