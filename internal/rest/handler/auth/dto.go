package authhandler

import (
	"time"

	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"user_name" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"user_name,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Password *string `json:"password,omitempty" binding:"omitempty,min=6"`
	Avatar   *string `json:"avatar,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UserDTO struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"user_name"`
	Name         string    `json:"name"`
	Avatar       *string   `json:"avatar,omitempty"`
	AuthProvider string    `json:"authProvider"`
	CreatedAt    time.Time `json:"createdAt"`
}

type AuthResponseData struct {
	Message      string  `json:"message,omitempty"`
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"accessToken"`
	RefreshToken string  `json:"refreshToken"`
}

func ToUserDTO(u *domain.User) UserDTO {
	return UserDTO{
		ID:           u.ID,
		Email:        u.Email,
		Username:     u.Username,
		Name:         u.Name,
		Avatar:       u.Avatar,
		AuthProvider: u.AuthProvider,
		CreatedAt:    u.CreatedAt,
	}
}
