package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/infrastructure/db/postgres"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/repository"
	userhandler "github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/handler/user"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/middleware"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/user"
	"go.uber.org/zap"
)

func runHTTPServer(ctx context.Context, cfg *config.Config, log *zap.Logger) error {
	gin.SetMode(cfg.Server.Mode)

	// Connect PostgreSQL
	conn, err := postgres.NewConnection(ctx, cfg.Database.ConnectionString())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close()
	log.Info("Successfully connected to the database")

	// Repositories
	userRepo := repository.NewUserRepository(conn.DB, log)

	// Domain Services
	userSvc := user.NewService(userRepo, log)

	// REST Handlers
	userHandler := userhandler.NewHandler(userSvc, log)

	// Router setup
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(100, 5*time.Minute))

	// API Routing
	aPIVersionOne(router, userHandler, log)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		log.Info("Server listening", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("ListenAndServe error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	return nil
}
