package domain

type Expense struct {
	ID       int      `json:"id"`
	Amount   int      `json:"amount"`
	Memo     string   `json:"memo"`
	SpentAt  string   `json:"spent_at"`
	Status   string   `json:"status"`
	Category Category `json:"category"`
}
