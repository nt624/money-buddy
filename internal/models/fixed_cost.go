package models

type FixedCost struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
