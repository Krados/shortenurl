package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	rdb redis.UniversalClient
}

func NewRedisRepository(rdb redis.UniversalClient) Repository {
	return &redisRepository{
		rdb: rdb,
	}
}

func (r *redisRepository) Get(key string) (res string, err error) {
	val, err := r.rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			err = ErrKeyNotFound
			return
		}
		return
	}
	res = val

	return
}

func (r *redisRepository) Set(key string, val string, expiration time.Duration) error {
	return r.rdb.Set(context.Background(), key, val, expiration).Err()
}
