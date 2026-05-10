package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/DimKa163/goseller/internal/configuration"
	"github.com/DimKa163/goseller/internal/middleware"
	"github.com/DimKa163/goseller/internal/shared/rest"
	"github.com/DimKa163/goseller/internal/user/infrastructure"
	userinterfaces "github.com/DimKa163/goseller/internal/user/interface"
	"github.com/DimKa163/goseller/internal/user/usecase"
	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	var userConfig configuration.UserConfiguration
	env.Parse(&userConfig)
	fmt.Printf("Loaded configuration: %+v\n", userConfig)
	log, err := createLogger(userConfig)
	if err != nil {
		panic(err)
	}
	db, err := createDatabaseConnection(ctx, userConfig.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()
	if err := infrastructure.Migrate(db, userConfig.MigrationPath); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}
	router := createServer()
	router.Use(middleware.LoggingMiddleware(log))
	healthCheck(router, db)
	mapControllers(router, userinterfaces.NewUserController(log, usecase.NewUser(infrastructure.NewUserRepository(db), log)))
	server := &http.Server{
		Addr:    userConfig.Addr,
		Handler: router,
	}
	go func() {
		<-ctx.Done()
		log.Info("Shutting down server...")
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Failed to shutdown server", zap.Error(err))
		}
		log.Info("Server gracefully stopped")
	}()
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server", zap.Error(err))
			return
		}
	}
}

func createLogger(conf configuration.UserConfiguration) (*zap.Logger, error) {
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
