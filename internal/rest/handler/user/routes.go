package userhandler

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(r *gin.RouterGroup, h *Handler) {
	users := r.Group("/users")
	{
		users.POST("/register", h.Register)
		users.GET("/:id", h.GetByID)
	}
}

func RegisterProtectedRoutes(r *gin.RouterGroup, h *Handler) {
	users := r.Group("/users")
	{
		users.GET("/me", h.Me)
	}
}
