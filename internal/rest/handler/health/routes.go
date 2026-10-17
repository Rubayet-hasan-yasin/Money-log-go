package healthhandler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	r.GET("/health", h.Health)
	r.GET("/health/detailed", h.DetailedHealth)
}
