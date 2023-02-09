package service

import (
	"github.com/Krados/shortenurl/internal/service/cache"
	"github.com/Krados/shortenurl/internal/service/shortener"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	shortener.New,
	shortener.NewMySQLRepository,
	cache.New,
	cache.NewRedisRepository,
)
