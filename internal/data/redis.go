package data

import (
	"github.com/Krados/shortenurl/internal/conf"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *conf.Config) redis.UniversalClient {
	return redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: cfg.Data.Redis.Addrs,
	})
}
