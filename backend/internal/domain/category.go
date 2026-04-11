package domain

type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CategoryType string `json:"category_type"`
	SortOrder    int    `json:"sort_order"`
}
