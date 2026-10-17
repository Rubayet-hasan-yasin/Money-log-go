package auth

import (
	"context"
	"strings"
	"testing"

	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
	"go.uber.org/zap"
)

type mockUserRepo struct {
	usersByID       map[string]*domain.User
	usersByEmail    map[string]*domain.User
	usersByUsername map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		usersByID:       make(map[string]*domain.User),
		usersByEmail:    make(map[string]*domain.User),
		usersByUsername: make(map[string]*domain.User),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, u *domain.User) error {
	m.usersByID[u.ID] = u
	m.usersByEmail[strings.ToLower(u.Email)] = u
	m.usersByUsername[strings.ToLower(u.Username)] = u
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if u, ok := m.usersByID[id]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if u, ok := m.usersByEmail[strings.ToLower(email)]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	if u, ok := m.usersByUsername[strings.ToLower(username)]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	for _, u := range m.usersByID {
		if u.GoogleID != nil && *u.GoogleID == googleID {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) Update(ctx context.Context, u *domain.User) error {
	m.usersByID[u.ID] = u
	m.usersByEmail[strings.ToLower(u.Email)] = u
	m.usersByUsername[strings.ToLower(u.Username)] = u
	return nil
}

func TestRegisterAndLoginWithUsername(t *testing.T) {
	repo := newMockUserRepo()
	logger := zap.NewNop()
	svc := NewService(repo, nil, nil, logger)

	ctx := context.Background()

	// 1. Register with unique username and separate display name
	res1, err := svc.Register(ctx, RegisterInput{
		Email:    "john@example.com",
		Username: "johndoe",
		Name:     "John Doe",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error registering user: %v", err)
	}

	if res1.User.Username != "johndoe" {
		t.Errorf("expected username 'johndoe', got '%s'", res1.User.Username)
	}
	if res1.User.Name != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", res1.User.Name)
	}

	// 2. Register duplicate username should fail with ErrUsernameAlreadyInUse
	_, err = svc.Register(ctx, RegisterInput{
		Email:    "another@example.com",
		Username: "johndoe",
		Name:     "Another John",
		Password: "password123",
	})
	if err != domain.ErrUsernameAlreadyInUse {
		t.Errorf("expected ErrUsernameAlreadyInUse, got %v", err)
	}

	// 3. Register duplicate email should fail with ErrEmailAlreadyInUse
	_, err = svc.Register(ctx, RegisterInput{
		Email:    "john@example.com",
		Username: "differentuser",
		Name:     "Different User",
		Password: "password123",
	})
	if err != domain.ErrEmailAlreadyInUse {
		t.Errorf("expected ErrEmailAlreadyInUse, got %v", err)
	}

	// 4. Register without username should fail with ErrInvalidInput
	_, err = svc.Register(ctx, RegisterInput{
		Email:    "jane@example.com",
		Name:     "Jane Smith",
		Password: "password123",
	})
	if err != domain.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput when username is missing, got: %v", err)
	}

	// 5. Login using email
	loginRes1, err := svc.Login(ctx, LoginInput{
		Email:    "john@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login with email failed: %v", err)
	}
	if loginRes1.User.ID != res1.User.ID {
		t.Errorf("expected user ID %s, got %s", res1.User.ID, loginRes1.User.ID)
	}

	// 6. Login using non-email should fail
	_, err = svc.Login(ctx, LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password123",
	})
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials for non-existent email, got %v", err)
	}
}

