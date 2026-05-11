package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/DimKa163/goseller/internal/configuration"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var config configuration.CategoryConfiguration
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	fmt.Printf("load conf %+v\n", config)
	log, err := createLogger(config)
	if err != nil {
		panic(err)
	}
	db, err := createDatabaseConnection(ctx, config.Database)
	if err != nil {
		log.Fatal("failed to connect", zap.Error(err))
	}
	defer db.Close()
}

func createLogger(conf configuration.CategoryConfiguration) (*zap.Logger, error) {
	if conf.AppEnv == "production" {
		return zap.NewProduction()
	} else {
		return zap.NewDevelopment()
	}
}

func createDatabaseConnection(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}
