package repository

import (
	"context"
	"database/sql"

	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/user"
	"go.uber.org/zap"
)

type UserRepository struct {
	db  *sql.DB
	log *zap.Logger
}

var _ user.Repository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB, log *zap.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `
	INSERT INTO users (id, email, username, password_hash, name, google_id, auth_provider, avatar, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.Username, u.Password, u.Name, u.GoogleID, u.AuthProvider, u.Avatar, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return domain.WrapErrorf(err, "failed to create user")
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
	SELECT id, email, username, password_hash, name, google_id, auth_provider, avatar, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.Name,
		&u.GoogleID,
		&u.AuthProvider,
		&u.Avatar,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, domain.WrapErrorf(err, "failed to find user by id")
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, username, password_hash, name, google_id, auth_provider, avatar, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.Name,
		&u.GoogleID,
		&u.AuthProvider,
		&u.Avatar,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, domain.WrapErrorf(err, "failed to find user by email")
	}
	return &u, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
	SELECT id, email, username, password_hash, name, google_id, auth_provider, avatar, created_at, updated_at
	FROM users
	WHERE LOWER(username) = LOWER($1)
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.Name,
		&u.GoogleID,
		&u.AuthProvider,
		&u.Avatar,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, domain.WrapErrorf(err, "failed to find user by username")
	}
	return &u, nil
}

func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*domain.User, error) {
	query := `
	SELECT id, email, username, password_hash, name, google_id, auth_provider, avatar, created_at, updated_at
	FROM users
	WHERE google_id = $1
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, googleID).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.Name,
		&u.GoogleID,
		&u.AuthProvider,
		&u.Avatar,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, domain.WrapErrorf(err, "failed to find user by google id")
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *domain.User) error {
	query := `
	UPDATE users
	SET email = $2, username = $3, password_hash = $4, name = $5, google_id = $6, auth_provider = $7, avatar = $8, updated_at = $9
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.Username, u.Password, u.Name, u.GoogleID, u.AuthProvider, u.Avatar, u.UpdatedAt)
	if err != nil {
		return domain.WrapErrorf(err, "failed to update user")
	}
	return nil
}

