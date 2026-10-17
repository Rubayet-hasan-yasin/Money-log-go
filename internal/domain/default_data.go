package domain

type DefaultCategory struct {
	Name  string
	Color string
	Icon  string
	Type  CategoryType
}

type DefaultWallet struct {
	Name    string
	Color   string
	Icon    string
	Balance float64
}

var DefaultCategories = []DefaultCategory{
	{Name: "No Category", Color: "#9CA3AF", Icon: "📦", Type: CategoryTypeExpense},
	{Name: "Food & Dining", Color: "#EF4444", Icon: "🍔", Type: CategoryTypeExpense},
	{Name: "Transportation", Color: "#3B82F6", Icon: "🚗", Type: CategoryTypeExpense},
	{Name: "Shopping", Color: "#8B5CF6", Icon: "🛍️", Type: CategoryTypeExpense},
	{Name: "Entertainment", Color: "#EC4899", Icon: "🎬", Type: CategoryTypeExpense},
	{Name: "Bills & Utilities", Color: "#F59E0B", Icon: "💡", Type: CategoryTypeExpense},
	{Name: "Healthcare", Color: "#10B981", Icon: "🏥", Type: CategoryTypeExpense},
	{Name: "Education", Color: "#06B6D4", Icon: "📚", Type: CategoryTypeExpense},
	{Name: "Travel", Color: "#6366F1", Icon: "✈️", Type: CategoryTypeExpense},
	{Name: "Personal Care", Color: "#EC4899", Icon: "💇", Type: CategoryTypeExpense},
	{Name: "Groceries", Color: "#84CC16", Icon: "🛒", Type: CategoryTypeExpense},
	{Name: "Other", Color: "#64748B", Icon: "📌", Type: CategoryTypeExpense},

	// Income categories
	{Name: "Salary", Color: "#10B981", Icon: "💰", Type: CategoryTypeIncome},
	{Name: "Gifts", Color: "#EC4899", Icon: "🎁", Type: CategoryTypeIncome},
	{Name: "Pocket Money", Color: "#F59E0B", Icon: "💵", Type: CategoryTypeIncome},
	{Name: "Other Income", Color: "#64748B", Icon: "💳", Type: CategoryTypeIncome},
}

var DefaultWallets = []DefaultWallet{
	{Name: "Cash", Color: "#10B981", Icon: "💵", Balance: 0.0},
}
