package domain

import "time"

type CategoryType string

const (
	CategoryTypeExpense CategoryType = "EXPENSE"
	CategoryTypeIncome  CategoryType = "INCOME"
)

type Category struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Color     *string      `json:"color,omitempty"`
	Icon      *string      `json:"icon,omitempty"`
	Type      CategoryType `json:"type"` // EXPENSE or INCOME
	UserID    string       `json:"user_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
