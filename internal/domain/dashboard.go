package domain

type DashboardSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	NetBalance   float64 `json:"net_balance"`
}

type CategoryAnalytics struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Color        *string `json:"color,omitempty"`
	Icon         *string `json:"icon,omitempty"`
	TotalAmount  float64 `json:"total_amount"`
	Percentage   float64 `json:"percentage"`
	Count        int64   `json:"count"`
}

type MonthlyTrend struct {
	Month   string  `json:"month"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}
