package user

import (
	"context"

	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
)

// Service defines the business logic methods for User domain
type Service interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// Repository defines data access methods for User domain
type Repository interface {
	Create(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}
