package domain

type MonthlySettings struct {
	ID         int    `json:"id"`
	UserID     string `json:"user_id"`
	Year       int    `json:"year"`
	Month      int    `json:"month"`
	Income     int    `json:"income"`
	SavingGoal int    `json:"saving_goal"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
