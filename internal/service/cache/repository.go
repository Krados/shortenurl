package cache

import "time"

//go:generate mockery --name=Repository --filename=mock_repository.go --inpackage
type Repository interface {
	Get(key string) (res string, err error)
	Set(key string, val string, expiration time.Duration) error
}
