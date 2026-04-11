package domain

type MonthlySettings struct {
	ID         int
	UserID     string
	Year       int
	Month      int
	Income     int
	SavingGoal int
	CreatedAt  string
	UpdatedAt  string
}
