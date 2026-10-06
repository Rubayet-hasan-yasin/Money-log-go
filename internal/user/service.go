package user

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) Service {
	return &service{
		repo: repo,
		log:  log,
	}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	username := strings.TrimSpace(input.Username)

	if email == "" || username == "" || len(input.Password) < 6 {
		return nil, domain.ErrInvalidInput
	}

	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, domain.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	now := time.Now()
	u := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		s.log.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	return u, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.FindByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.FindByEmail(ctx, strings.ToLower(email))
}
