package usecase

import (
	"context"

	"money-buddy-backend/internal/domain"
)

type getUserRepository interface {
	GetUserByID(ctx context.Context, id string) (domain.User, error)
}

// GetUserUseCase はユーザー取得ユースケースです。
type GetUserUseCase struct {
	repo getUserRepository
}

// NewGetUserUseCase は GetUserUseCase の新しいインスタンスを生成します。
func NewGetUserUseCase(repo getUserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

// Execute はユーザー取得ユースケースを実行します。
func (uc *GetUserUseCase) Execute(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
