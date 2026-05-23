package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/DimKa163/goseller/internal/category/infrastructure"
	categoryinterfaces "github.com/DimKa163/goseller/internal/category/interfaces"
	"github.com/DimKa163/goseller/internal/category/usecase"
	"github.com/DimKa163/goseller/internal/configuration"
	"github.com/DimKa163/goseller/internal/middleware"
	"github.com/DimKa163/goseller/internal/shared/rest"
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
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
	if err := infrastructure.Migrate(db, config.MigrationPath); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}
	router := createServer()
	router.Use(middleware.LoggingMiddleware(log))
	healthCheck(router, db)
	mapControllers(router, categoryinterfaces.NewCategoryController(log, usecase.NewCategoryAppService(log, infrastructure.NewCategoryRepository(db))))
	server := &http.Server{
		Addr:    config.Addr,
		Handler: router,
	}
	go func() {
		<-ctx.Done()
		log.Info("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Failed to shutdown server", zap.Error(err))
		}
		log.Info("Server gracefully stopped")
	}()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server", zap.Error(err))
			return
		}
	}
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

func createServer() *gin.Engine {
	router := gin.New()
	return router
}

func healthCheck(router *gin.Engine, db *pgxpool.Pool) {
	router.GET("/healthCheck", func(ctx *gin.Context) {
		if err := db.Ping(ctx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "OK", "database": "Connected"})
	})
}

func mapControllers(router *gin.Engine, controllers ...rest.Controller) {
	for _, controller := range controllers {
		controller.Map(router)
	}
}
