package domain

import "errors"

// ErrDuplicateCategoryName は同一ユーザー内でカテゴリ名が重複している場合のエラーです。
var ErrDuplicateCategoryName = errors.New("category name already exists")

// ErrCategoryNotFound はカテゴリが存在しない場合のエラーです。
var ErrCategoryNotFound = errors.New("category not found")

type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CategoryType string `json:"category_type"`
	SortOrder    int    `json:"sort_order"`
}
