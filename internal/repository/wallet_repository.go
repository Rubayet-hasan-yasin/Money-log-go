package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/domain"
	"go.uber.org/zap"
)

type WalletRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewWalletRepository(db *sql.DB, log *zap.Logger) *WalletRepository {
	return &WalletRepository{
		db:  db,
		log: log,
	}
}

func (r *WalletRepository) CreateBatch(ctx context.Context, wallets []domain.Wallet) error {
	if len(wallets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.WrapErrorf(err, "failed to begin transaction")
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO wallets (id, name, color, icon, balance, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, name) DO NOTHING
	`)
	if err != nil {
		return domain.WrapErrorf(err, "failed to prepare wallet statement")
	}
	defer stmt.Close()

	for _, w := range wallets {
		id := w.ID
		if id == "" {
			id = uuid.New().String()
		}
		now := time.Now()
		createdAt := w.CreatedAt
		if createdAt.IsZero() {
			createdAt = now
		}
		updatedAt := w.UpdatedAt
		if updatedAt.IsZero() {
			updatedAt = now
		}

		_, err := stmt.ExecContext(ctx, id, w.Name, w.Color, w.Icon, w.Balance, w.UserID, createdAt, updatedAt)
		if err != nil {
			return domain.WrapErrorf(err, "failed to insert wallet: %s", w.Name)
		}
	}

	return tx.Commit()
}

func (r *WalletRepository) CreateDefaultWallets(ctx context.Context, userID string) error {
	var wallets []domain.Wallet
	now := time.Now()
	for _, dw := range domain.DefaultWallets {
		color := dw.Color
		icon := dw.Icon
		wallets = append(wallets, domain.Wallet{
			ID:        uuid.New().String(),
			Name:      dw.Name,
			Color:     &color,
			Icon:      &icon,
			Balance:   dw.Balance,
			UserID:    userID,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return r.CreateBatch(ctx, wallets)
}
