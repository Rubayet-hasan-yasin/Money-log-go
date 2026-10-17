package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/auth"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/infrastructure/db/postgres"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/repository"
	authhandler "github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/handler/auth"
	healthhandler "github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/handler/health"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/middleware"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

	// Configure Google OAuth2 provider
	var googleOAuthConfig *oauth2.Config
	if cfg.OAuth.GoogleClientID != "" && cfg.OAuth.GoogleClientSecret != "" && cfg.OAuth.GoogleCallbackURL != "" {
		googleOAuthConfig = &oauth2.Config{
			ClientID:     cfg.OAuth.GoogleClientID,
			ClientSecret: cfg.OAuth.GoogleClientSecret,
			RedirectURL:  cfg.OAuth.GoogleCallbackURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
		log.Info("Google OAuth provider configured")
	} else {
		log.Warn("Google OAuth credentials not configured; OAuth login disabled")
	}

	// Repositories
	userRepo := repository.NewUserRepository(conn.DB, log)
	categoryRepo := repository.NewCategoryRepository(conn.DB, log)
	walletRepo := repository.NewWalletRepository(conn.DB, log)

	// Domain Services
	authSvc := auth.NewService(userRepo, categoryRepo, walletRepo, log)

	// REST Handlers
	healthHandler := healthhandler.NewHandler(conn.DB, log)
	authHandler := authhandler.NewHandler(authSvc, log, googleOAuthConfig)

	// Router setup
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(100, 5*time.Minute))

	// API Routing
	aPIVersionOne(router, healthHandler, authHandler, log)

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
