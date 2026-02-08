package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"money-buddy-backend/internal/models"
)

// mockFixedCostRepo はテスト用のモックリポジトリです
type mockFixedCostRepo struct {
	mock.Mock
}

func (m *mockFixedCostRepo) CreateFixedCost(ctx context.Context, userID string, name string, amount int) (models.FixedCost, error) {
	args := m.Called(ctx, userID, name, amount)
	return args.Get(0).(models.FixedCost), args.Error(1)
}

func (m *mockFixedCostRepo) ListFixedCostsByUser(ctx context.Context, userID string) ([]models.FixedCost, error) {
	args := m.Called(ctx, userID)
	if result := args.Get(0); result != nil {
		return result.([]models.FixedCost), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockFixedCostRepo) DeleteFixedCostsByUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockFixedCostRepo) BulkCreateFixedCosts(ctx context.Context, userID string, fixedCosts []models.FixedCostInput) error {
	args := m.Called(ctx, userID, fixedCosts)
	return args.Error(0)
}

func (m *mockFixedCostRepo) UpdateFixedCost(ctx context.Context, id int32, userID string, name string, amount int) error {
	args := m.Called(ctx, id, userID, name, amount)
	return args.Error(0)
}

func (m *mockFixedCostRepo) DeleteFixedCost(ctx context.Context, id int32, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

// TestListFixedCosts は固定費一覧取得のテストです
func TestListFixedCosts(t *testing.T) {
	ctx := context.Background()

	t.Run("正常に固定費一覧を取得できる", func(t *testing.T) {
		repo := new(mockFixedCostRepo)
		expected := []models.FixedCost{
			{ID: 1, UserID: "user1", Name: "家賃", Amount: 80000},
			{ID: 2, UserID: "user1", Name: "通信費", Amount: 5000},
		}
		repo.On("ListFixedCostsByUser", ctx, "user1").Return(expected, nil)

		service := NewFixedCostService(repo)
		result, err := service.ListFixedCosts(ctx, "user1")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		repo.AssertExpectations(t)
	})

	t.Run("空の一覧を取得できる", func(t *testing.T) {
		repo := new(mockFixedCostRepo)
		repo.On("ListFixedCostsByUser", ctx, "user1").Return([]models.FixedCost{}, nil)

		service := NewFixedCostService(repo)
		result, err := service.ListFixedCosts(ctx, "user1")

		assert.NoError(t, err)
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})
}

// TestUpdateFixedCost は固定費更新のテストです
func TestUpdateFixedCost(t *testing.T) {
	ctx := context.Background()

	t.Run("正常に固定費を更新できる", func(t *testing.T) {
		repo := new(mockFixedCostRepo)
		repo.On("UpdateFixedCost", ctx, int32(1), "user1", "家賃（更新）", 85000).Return(nil)
		repo.On("ListFixedCostsByUser", ctx, "user1").Return([]models.FixedCost{
			{ID: 1, UserID: "user1", Name: "家賃（更新）", Amount: 85000},
		}, nil)

		service := NewFixedCostService(repo)
		result, err := service.UpdateFixedCost(ctx, "user1", 1, "家賃（更新）", 85000)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "家賃（更新）", result.Name)
		assert.Equal(t, 85000, result.Amount)
		repo.AssertExpectations(t)
	})

	t.Run("名前が空の場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, "", 80000)

		assert.Error(t, err)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("名前が空白のみの場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, "   ", 80000)

		assert.Error(t, err)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve)
	})

	t.Run("名前が最大長を超える場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)
		longName := string(make([]byte, 101)) // 101文字

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, longName, 80000)

		assert.Error(t, err)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve)
		assert.Contains(t, err.Error(), "name is too long")
	})

	t.Run("金額が0以下の場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, "家賃", 0)

		assert.Error(t, err)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve)
		assert.Contains(t, err.Error(), "amount must be greater than 0")
	})

	t.Run("金額が負の値の場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, "家賃", -1000)

		assert.Error(t, err)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve)
	})

	t.Run("金額が上限を超える場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, "家賃", 1000000001)

		assert.Error(t, err)
		var ve *ValidationError
		assert.ErrorAs(t, err, &ve)
		assert.Contains(t, err.Error(), "amount exceeds maximum allowed")
	})

	t.Run("更新後に固定費が見つからない場合はエラー", func(t *testing.T) {
		repo := new(mockFixedCostRepo)
		repo.On("UpdateFixedCost", ctx, int32(1), "user1", "家賃", 80000).Return(nil)
		repo.On("ListFixedCostsByUser", ctx, "user1").Return([]models.FixedCost{}, nil)

		service := NewFixedCostService(repo)
		_, err := service.UpdateFixedCost(ctx, "user1", 1, "家賃", 80000)

		assert.Error(t, err)
		var ne *NotFoundError
		assert.ErrorAs(t, err, &ne)
		repo.AssertExpectations(t)
	})
}

// TestDeleteFixedCost は固定費削除のテストです
func TestDeleteFixedCost(t *testing.T) {
	ctx := context.Background()

	t.Run("正常に固定費を削除できる", func(t *testing.T) {
		repo := new(mockFixedCostRepo)
		repo.On("DeleteFixedCost", ctx, int32(1), "user1").Return(nil)

		service := NewFixedCostService(repo)
		err := service.DeleteFixedCost(ctx, "user1", 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
