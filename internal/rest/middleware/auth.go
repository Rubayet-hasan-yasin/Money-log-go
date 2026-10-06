package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
	"go.uber.org/zap"
)

// AuthMiddleware verifies the Bearer JWT token from the Authorization header
func AuthMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be Bearer token"})
			return
		}

		tokenString := parts[1]
		jwtSecret := config.GetConfig().JWT.Secret

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Warn("invalid auth token", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, exists := claims["user_id"]; exists {
				c.Set("user_id", fmt.Sprint(userID))
			}
		}

		c.Next()
	}
}

// GetUserID retrieves the authenticated user ID from context
func GetUserID(c *gin.Context) string {
	if val, ok := c.Get("user_id"); ok {
		return val.(string)
	}
	return ""
}
