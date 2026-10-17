package domain

import "time"

type TransactionType string

const (
	TransactionTypeExpense  TransactionType = "EXPENSE"
	TransactionTypeIncome   TransactionType = "INCOME"
	TransactionTypeTransfer TransactionType = "TRANSFER"

	DefaultCurrency = "BDT"
)

type Expense struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Amount      float64         `json:"amount"`
	Currency    string          `json:"currency"`
	Type        TransactionType `json:"type"` // EXPENSE, INCOME, TRANSFER
	CategoryID  *string         `json:"category_id,omitempty"`
	Category    *Category       `json:"category,omitempty"`
	WalletID    *string         `json:"wallet_id,omitempty"`
	Wallet      *Wallet         `json:"wallet,omitempty"`
	ToWalletID  *string         `json:"to_wallet_id,omitempty"`
	ToWallet    *Wallet         `json:"to_wallet,omitempty"`
	Description *string         `json:"description,omitempty"`
	Date        time.Time       `json:"date"`
	UserID      string          `json:"user_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
