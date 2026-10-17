package main

import (
	"github.com/gin-gonic/gin"
	authhandler "github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/handler/auth"
	healthhandler "github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/handler/health"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/middleware"
	"go.uber.org/zap"
)

func aPIVersionOne(
	router *gin.Engine,
	healthHandler *healthhandler.Handler,
	authHandler *authhandler.Handler,
	log *zap.Logger,
) {
	api := router.Group("/api/v1")
	{
		// Health check endpoints
		healthhandler.RegisterRoutes(api, healthHandler)

		// Public route groups
		authhandler.RegisterPublicRoutes(api, authHandler)

		// Protected route groups (authenticated via Bearer JWT)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(log))
		{
			authhandler.RegisterProtectedRoutes(protected, authHandler)
		}
	}
}
