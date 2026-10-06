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
	INSERT INTO users (id, email, user_name, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, u.ID, u.Email, u.Username, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return domain.WrapErrorf(err, "failed to create user")
	}
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
	SELECT id, email, user_name, password_hash, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
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
	SELECT id, email, user_name, password_hash, created_at, updated_at
	FROM users
	WHERE email = $1
	`
	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Username,
		&u.Password,
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
