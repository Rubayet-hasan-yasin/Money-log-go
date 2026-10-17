package authhandler

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(r *gin.RouterGroup, h *Handler) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.GET("/google", h.GoogleLogin)
		auth.GET("/google/callback", h.GoogleCallback)
		auth.GET("/google/failure", h.GoogleFailure)
	}
}

func RegisterProtectedRoutes(r *gin.RouterGroup, h *Handler) {
	auth := r.Group("/auth")
	{
		auth.GET("/profile", h.GetProfile)
		auth.PUT("/profile", h.UpdateProfile)
	}
}
