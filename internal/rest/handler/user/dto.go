package userhandler

import "github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"

type RegisterRequest struct {
	Username string `json:"user_name" binding:"required,min=3,max=50"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"user_name"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func ToUserResponse(u *domain.User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
