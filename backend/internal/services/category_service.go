package services

import (
	"context"

	"money-buddy-backend/internal/models"
)

type CategoryService interface {
	ListCategories(ctx context.Context) ([]models.Category, error)
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) ListCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListCategories(ctx)
}
