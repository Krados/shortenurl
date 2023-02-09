package cache

import (
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type cache struct {
	repo Repository
}

func New(repo Repository) Cache {
	return &cache{
		repo: repo,
	}
}

func (c *cache) Get(key string) (res string, err error) {
	return c.repo.Get(key)
}

func (c *cache) Set(key string, val string, expiration time.Duration) error {
	return c.repo.Set(key, val, expiration)
}
