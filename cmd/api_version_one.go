package main

import (
	"github.com/gin-gonic/gin"
	userhandler "github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/handler/user"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/middleware"
	"go.uber.org/zap"
)

func aPIVersionOne(
	router *gin.Engine,
	userHandler *userhandler.Handler,
	log *zap.Logger,
) {
	api := router.Group("/api/v1")
	{
		// Health check endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "go-backend-skeleton",
			})
		})

		// Public route group
		userhandler.RegisterPublicRoutes(api, userHandler)

		// Protected route group (authenticated via Bearer JWT)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(log))
		{
			userhandler.RegisterProtectedRoutes(protected, userHandler)
		}
	}
}
