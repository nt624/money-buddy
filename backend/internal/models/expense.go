package models

type CreateExpenseInput struct {
	Amount     *int   `json:"amount"`
	CategoryID *int   `json:"category_id"`
	Memo       string `json:"memo"`
	SpentAt    string `json:"spent_at"`
	Status     string `json:"status"`
}

type UpdateExpenseInput struct {
	ID         int    `json:"id"`
	Amount     *int   `json:"amount"`
	CategoryID *int   `json:"category_id"`
	Memo       string `json:"memo"`
	SpentAt    string `json:"spent_at"`
	Status     string `json:"status"`
}

type Expense struct {
	ID       int      `json:"id"`
	Amount   int      `json:"amount"`
	Memo     string   `json:"memo"`
	SpentAt  string   `json:"spent_at"`
	Status   string   `json:"status"`
	Category Category `json:"category"`
}
