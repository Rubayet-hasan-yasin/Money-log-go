package auth

import (
	"context"

	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
)

type RegisterInput struct {
	Email    string
	Username string
	Password string
	Name     string
}

type LoginInput struct {
	Email    string
	Password string
}

type UpdateProfileInput struct {
	Email    *string
	Username *string
	Password *string
	Name     *string
	Avatar   *string
}

type OAuthProfile struct {
	Provider string
	Email    string
	Name     string
	Avatar   string
	GoogleID string
}

type AuthResult struct {
	User         *domain.User
	AccessToken  string // Access token
	RefreshToken string // Refresh token
	IsNewUser    bool
}

type TokenClaims struct {
	UserID    string
	Email     string
	Username  string
	Name      string
	TokenType string
}

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
}

type CategoryRepo interface {
	CreateDefaultCategories(ctx context.Context, userID string) error
}

type WalletRepo interface {
	CreateDefaultWallets(ctx context.Context, userID string) error
}

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*domain.User, error)
	OAuthLogin(ctx context.Context, profile OAuthProfile) (*AuthResult, error)
	ValidateToken(ctx context.Context, tokenString string) (*TokenClaims, error)
	RefreshToken(ctx context.Context, tokenString string) (*AuthResult, error)
}
