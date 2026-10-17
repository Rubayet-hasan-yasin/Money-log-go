package userhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/middleware"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/rest/response"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/user"
	"go.uber.org/zap"
)

type Handler struct {
	svc user.Service
	log *zap.Logger
}

func NewHandler(svc user.Service, log *zap.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	u, err := h.svc.Register(c.Request.Context(), user.RegisterInput{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, ToUserResponse(u))
}

func (h *Handler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		response.Unauthorized(c, "unauthorized")
		return
	}

	u, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, ToUserResponse(u))
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	u, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, ToUserResponse(u))
}
