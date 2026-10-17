package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
	"go.uber.org/zap"
)

type CategoryRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewCategoryRepository(db *sql.DB, log *zap.Logger) *CategoryRepository {
	return &CategoryRepository{
		db:  db,
		log: log,
	}
}

func (r *CategoryRepository) CreateBatch(ctx context.Context, categories []domain.Category) error {
	if len(categories) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.WrapErrorf(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO categories (id, name, color, icon, type, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, name) DO NOTHING
	`)
	if err != nil {
		return domain.WrapErrorf(err, "failed to prepare category statement")
	}
	defer stmt.Close()

	for _, c := range categories {
		id := c.ID
		if id == "" {
			id = uuid.New().String()
		}
		now := time.Now()
		createdAt := c.CreatedAt
		if createdAt.IsZero() {
			createdAt = now
		}
		updatedAt := c.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = now
		}

		_, err := stmt.ExecContext(ctx, id, c.Name, c.Color, c.Icon, c.Type, c.UserID, createdAt, updatedAt)
		if err != nil {
			return domain.WrapErrorf(err, "failed to insert category: %s", c.Name)
		}
	}

	return tx.Commit()
}

func (r *CategoryRepository) CreateDefaultCategories(ctx context.Context, userID string) error {
	var categories []domain.Category
	now := time.Now()
	for _, dc := range domain.DefaultCategories {
		color := dc.Color
		icon := dc.Icon
		categories = append(categories, domain.Category{
			ID:        uuid.New().String(),
			Name:      dc.Name,
			Color:     &color,
			Icon:      &icon,
			Type:      dc.Type,
			UserID:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return r.CreateBatch(ctx, categories)
}
