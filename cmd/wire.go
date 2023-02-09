//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Krados/shortenurl/internal/conf"
	"github.com/Krados/shortenurl/internal/controller"
	"github.com/Krados/shortenurl/internal/data"
	"github.com/Krados/shortenurl/internal/server"
	"github.com/Krados/shortenurl/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func initApp(cfg *conf.Config, sugar *zap.SugaredLogger) (*gin.Engine, error) {
	panic(wire.Build(controller.ProviderSet, data.ProviderSet, server.ProviderSet, service.ProviderSet))
}
