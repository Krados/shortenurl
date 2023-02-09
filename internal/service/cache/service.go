package cache

import "time"

//go:generate mockery --name=Cache --filename=mock_cache.go --inpackage
type Cache interface {
	Get(key string) (res string, err error)
	Set(key string, val string, expiration time.Duration) error
}
