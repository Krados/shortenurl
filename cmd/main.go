package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Krados/shortenurl/internal/conf"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// load config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(flagconf)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var cfg conf.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	// init logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// init app
	ginEngine, err := initApp(&cfg, sugar)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    cfg.Server.HTTP.Addr,
		Handler: ginEngine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		sugar.Fatal("Server Shutdown: ", err)
	}

	sugar.Infoln("Server exiting")
}
