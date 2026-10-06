package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
)

type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
}

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "success",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, message string, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "created successfully",
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}

func Error(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	switch err {
	case domain.ErrNotFound:
		statusCode = http.StatusNotFound
		message = err.Error()
	case domain.ErrAlreadyExists:
		statusCode = http.StatusConflict
		message = err.Error()
	case domain.ErrInvalidInput:
		statusCode = http.StatusBadRequest
		message = err.Error()
	case domain.ErrUnauthorized:
		statusCode = http.StatusUnauthorized
		message = err.Error()
	case domain.ErrForbidden:
		statusCode = http.StatusForbidden
		message = err.Error()
	default:
		if err != nil {
			message = err.Error()
		}
	}

	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
	})
}
