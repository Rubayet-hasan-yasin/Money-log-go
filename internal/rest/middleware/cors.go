package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	cfg "github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
)

func CORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	appConfig := cfg.GetConfig()

	if appConfig.Server.AllowedOrigins == "*" || appConfig.Server.AllowedOrigins == "" {
		corsConfig.AllowAllOrigins = true
	} else {
		var allowedOrigins []string
		origins := strings.Split(appConfig.Server.AllowedOrigins, ",")
		for _, origin := range origins {
			if trimmed := strings.TrimSpace(origin); trimmed != "" {
				allowedOrigins = append(allowedOrigins, trimmed)
			}
		}
		corsConfig.AllowOrigins = allowedOrigins
	}

	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
	corsConfig.AllowCredentials = true

	return cors.New(corsConfig)
}
