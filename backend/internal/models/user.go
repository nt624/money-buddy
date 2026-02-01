package models

type User struct {
	ID         string `json:"id"`
	Income     int    `json:"income"`
	SavingGoal int    `json:"saving_goal"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
