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
	username := strings.TrimSpace(strings.ToLower(input.Username))
	name := strings.TrimSpace(input.Name)
	if name == "" {
		name = username
	}

	if email == "" || username == "" || len(input.Password) < 6 {
		return nil, domain.ErrInvalidInput
	}

	existingEmail, _ := s.repo.FindByEmail(ctx, email)
	if existingEmail != nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	existingUsername, _ := s.repo.FindByUsername(ctx, username)
	if existingUsername != nil {
		return nil, domain.ErrUsernameAlreadyInUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, domain.ErrInternalServer
	}

	hashStr := string(hashedPassword)
	now := time.Now()
	u := &domain.User{
		ID:           uuid.New().String(),
		Email:        email,
		Username:     username,
		Password:     &hashStr,
		Name:         name,
		AuthProvider: domain.AuthProviderLocal,
		CreatedAt:    now,
		UpdatedAt:    now,
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

func (s *service) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	if strings.TrimSpace(username) == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.FindByUsername(ctx, strings.ToLower(username))
}

