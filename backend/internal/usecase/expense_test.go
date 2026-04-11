package usecase

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pace-wallet-backend/internal/domain"
)

// --- mocks for create expense ---

type mockCreateExpenseRepo struct {
	called bool
	in     CreateExpenseInput
}

func (m *mockCreateExpenseRepo) CreateExpense(userID string, input CreateExpenseInput) (domain.Expense, error) {
	m.called = true
	m.in = input
	return domain.Expense{ID: 1, Amount: *input.Amount, Memo: input.Memo, SpentAt: input.SpentAt, Category: domain.Category{ID: *input.CategoryID}}, nil
}

func (m *mockCreateExpenseRepo) GetExpenseByID(userID string, id int32) (domain.Expense, error) {
	return domain.Expense{}, errors.New("not implemented")
}

type mockCategoryRepo struct {
	exists map[int32]bool
	err    error
}

func (m *mockCategoryRepo) CategoryExists(ctx context.Context, userID string, id int32) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	if m.exists == nil {
		return false, nil
	}
	v, ok := m.exists[id]
	if !ok {
		return false, nil
	}
	return v, nil
}

func TestCreateExpenseValidation(t *testing.T) {
	cases := []struct {
		name       string
		input      CreateExpenseInput
		wantErr    bool
		wantCalled bool
	}{
		{name: "金額が0以下の場合はエラーになる", input: CreateExpenseInput{Amount: intPtr(0), CategoryID: intPtr(1), SpentAt: "2020-01-01"}, wantErr: true, wantCalled: false},
		{name: "金額が1Bを超える場合はエラーになる", input: CreateExpenseInput{Amount: intPtr(1000000001), CategoryID: intPtr(1), SpentAt: "2020-01-02"}, wantErr: true, wantCalled: false},
		{name: "カテゴリIDが0以下の場合はエラーになる", input: CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(0), SpentAt: "2020-01-01"}, wantErr: true, wantCalled: false},
		{name: "カテゴリが存在しない場合はエラーになる（カテゴリテーブル参照）", input: CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(9999), SpentAt: "2020-01-02"}, wantErr: true, wantCalled: false},
		{name: "spentAtがzero time（0001-01-01T00:00:00Z）の場合はエラーになる", input: CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: "0001-01-01T00:00:00Z"}, wantErr: true, wantCalled: false},
		{name: "日付のみ（YYYY-MM-DD）で正常に作成される", input: CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: "2020-01-02"}, wantErr: false, wantCalled: true},
		{name: "RFC3339形式で正常に作成される", input: CreateExpenseInput{Amount: intPtr(200), CategoryID: intPtr(2), SpentAt: "2020-01-02T15:04:05Z"}, wantErr: false, wantCalled: true},
		{name: "spentAtが空文字の場合はエラーになる", input: CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: ""}, wantErr: true, wantCalled: false},
		{name: "メモが最大長を超える場合はエラーになる", input: CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: "2020-01-02", Memo: strings.Repeat("a", 5001)}, wantErr: true, wantCalled: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := &mockCreateExpenseRepo{}
			exists := map[int32]bool{}
			if tc.input.CategoryID != nil && tc.wantCalled {
				exists[int32(*tc.input.CategoryID)] = true
			}
			cr := &mockCategoryRepo{exists: exists}
			uc := NewCreateExpenseUseCase(m, cr)

			out, err := uc.Execute("test-user", tc.input)

			if tc.wantErr {
				if !assert.Error(t, err, "expected error for case %s", tc.name) {
					return
				}
				var ve *ValidationError
				assert.ErrorAs(t, err, &ve, "expected ValidationError for case %s", tc.name)
			} else {
				assert.NoError(t, err, "unexpected error for case %s: %v", tc.name, err)
			}

			assert.Equal(t, tc.wantCalled, m.called, "repo called mismatch for case %s", tc.name)

			if !tc.wantErr {
				assert.Equal(t, *tc.input.Amount, out.Amount, "amount mismatch for case %s", tc.name)
			}
		})
	}
}

func intPtr(v int) *int { return &v }

func TestCreateExpense_DBErrorMapping(t *testing.T) {
	t.Parallel()

	validInput := CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: "2020-01-02"}

	cases := []struct {
		name     string
		repoErr  error
		wantType string
		wantMsg  string
	}{
		{name: "sql.ErrNoRows の場合は NotFoundError", repoErr: sql.ErrNoRows, wantType: "NotFoundError", wantMsg: "見つかりません"},
		{name: "外部キー(category_id)違反は ValidationError", repoErr: errors.New("pq: insert or update on table \"expenses\" violates foreign key constraint \"expenses_category_id_fkey\""), wantType: "ValidationError", wantMsg: "カテゴリが存在しません"},
		{name: "その他の DB エラーは InternalError", repoErr: errors.New("some db problem"), wantType: "InternalError", wantMsg: "internal error"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			m := &mockCreateExpenseRepoErr{returnErr: tc.repoErr}
			cr := &mockCategoryRepo{exists: map[int32]bool{1: true}}
			uc := NewCreateExpenseUseCase(m, cr)

			_, err := uc.Execute("test-user", validInput)
			if !assert.Error(t, err) {
				return
			}

			switch tc.wantType {
			case "NotFoundError":
				var e *NotFoundError
				if !assert.ErrorAs(t, err, &e) {
					return
				}
				assert.Contains(t, e.Error(), tc.wantMsg)
			case "ValidationError":
				var e *ValidationError
				if !assert.ErrorAs(t, err, &e) {
					return
				}
				assert.Equal(t, tc.wantMsg, e.Message)
			case "InternalError":
				var e *InternalError
				if !assert.ErrorAs(t, err, &e) {
					return
				}
				assert.Contains(t, e.Error(), tc.wantMsg)
			default:
				t.Fatalf("unknown wantType: %s", tc.wantType)
			}
		})
	}
}

type mockCreateExpenseRepoErr struct {
	returnErr error
}

func (m *mockCreateExpenseRepoErr) CreateExpense(userID string, input CreateExpenseInput) (domain.Expense, error) {
	return domain.Expense{}, m.returnErr
}

func (m *mockCreateExpenseRepoErr) GetExpenseByID(userID string, id int32) (domain.Expense, error) {
	return domain.Expense{}, errors.New("not implemented")
}

func TestCreateExpense_CategoryExistsError(t *testing.T) {
	t.Parallel()

	input := CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: "2020-01-02"}

	m := &mockCreateExpenseRepo{}
	cr := &mockCategoryRepo{err: errors.New("db error")}
	uc := NewCreateExpenseUseCase(m, cr)

	_, err := uc.Execute("test-user", input)
	if err == nil {
		t.Fatalf("expected error")
	}
	var ie *InternalError
	if !assert.ErrorAs(t, err, &ie) {
		t.Fatalf("expected InternalError, got %v", err)
	}
	assert.False(t, m.called, "repo should not be called when category existence check fails")
}

func TestCreateExpense_StatusValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		input          CreateExpenseInput
		wantErr        bool
		wantCalled     bool
		wantNormalized string
	}{
		{
			name:           "status 'planned' accepted and normalized",
			input:          CreateExpenseInput{Amount: intPtr(100), CategoryID: intPtr(1), SpentAt: "2020-01-02", Status: "PlAnNeD"},
			wantErr:        false,
			wantCalled:     true,
			wantNormalized: "planned",
		},
		{
			name:           "status 'confirmed' accepted",
			input:          CreateExpenseInput{Amount: intPtr(200), CategoryID: intPtr(2), SpentAt: "2020-01-03", Status: "confirmed"},
			wantErr:        false,
			wantCalled:     true,
			wantNormalized: "confirmed",
		},
		{
			name:       "invalid status causes ValidationError",
			input:      CreateExpenseInput{Amount: intPtr(300), CategoryID: intPtr(3), SpentAt: "2020-01-04", Status: "draft"},
			wantErr:    true,
			wantCalled: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &mockCreateExpenseRepo{}
			exists := map[int32]bool{}
			if tc.input.CategoryID != nil {
				exists[int32(*tc.input.CategoryID)] = true
			}
			cr := &mockCategoryRepo{exists: exists}
			uc := NewCreateExpenseUseCase(m, cr)

			_, err := uc.Execute("test-user", tc.input)

			if tc.wantErr {
				if !assert.Error(t, err) {
					return
				}
				var ve *ValidationError
				assert.ErrorAs(t, err, &ve)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantNormalized, m.in.Status)
			}

			assert.Equal(t, tc.wantCalled, m.called)
		})
	}
}

// --- mocks for delete expense ---

type mockDeleteExpenseRepo struct {
	called    bool
	deletedID int32
	returnErr error
}

func (m *mockDeleteExpenseRepo) GetExpenseByID(userID string, id int32) (domain.Expense, error) {
	if id == 9999 {
		return domain.Expense{}, sql.ErrNoRows
	}
	status := "confirmed"
	if id == 10 {
		status = "planned"
	}
	return domain.Expense{ID: int(id), Status: status}, nil
}

func (m *mockDeleteExpenseRepo) DeleteExpense(userID string, id int32) error {
	m.called = true
	m.deletedID = id
	return m.returnErr
}

func TestDeleteExpense_Success(t *testing.T) {
	t.Parallel()

	repo := &mockDeleteExpenseRepo{returnErr: nil}
	uc := NewDeleteExpenseUseCase(repo)

	err := uc.Execute("test-user", 1)
	assert.NoError(t, err)
	assert.True(t, repo.called)
	assert.Equal(t, int32(1), repo.deletedID)
}

func TestDeleteExpense_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockDeleteExpenseRepo{}
	uc := NewDeleteExpenseUseCase(repo)

	err := uc.Execute("test-user", 9999)
	var nfe *NotFoundError
	assert.ErrorAs(t, err, &nfe)
}

func TestDeleteExpense_StatusAgnostic(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		id   int
	}{
		{name: "planned item can delete", id: 10},
		{name: "confirmed item can delete", id: 11},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := &mockDeleteExpenseRepo{returnErr: nil}
			uc := NewDeleteExpenseUseCase(repo)

			err := uc.Execute("test-user", tc.id)
			assert.NoError(t, err)
			assert.True(t, repo.called)
			assert.Equal(t, int32(tc.id), repo.deletedID)
		})
	}
}

// --- mocks for update expense ---

type mockUpdateExpenseRepo struct {
	current   domain.Expense
	called    bool
	in        UpdateExpenseInput
	returnErr error
	getErr    error
}

func (m *mockUpdateExpenseRepo) GetExpenseByID(userID string, id int32) (domain.Expense, error) {
	if m.getErr != nil {
		return domain.Expense{}, m.getErr
	}
	return m.current, nil
}

func (m *mockUpdateExpenseRepo) UpdateExpense(userID string, input UpdateExpenseInput) (domain.Expense, error) {
	m.called = true
	m.in = input
	if m.returnErr != nil {
		return domain.Expense{}, m.returnErr
	}
	e := m.current
	if input.Amount != nil {
		e.Amount = *input.Amount
	}
	if input.CategoryID != nil {
		e.Category = domain.Category{ID: *input.CategoryID}
	}
	e.Memo = input.Memo
	e.SpentAt = input.SpentAt
	if input.Status != "" {
		e.Status = strings.ToLower(input.Status)
	}
	m.current = e
	return e, nil
}

func TestUpdateExpense_NormalCases(t *testing.T) {
	t.Parallel()

	t.Run("planned to confirmed", func(t *testing.T) {
		t.Parallel()

		repo := &mockUpdateExpenseRepo{current: domain.Expense{ID: 1, Amount: 100, Memo: "old", SpentAt: "2025-01-01", Status: "planned", Category: domain.Category{ID: 1}}}
		cr := &mockCategoryRepo{exists: map[int32]bool{2: true}}
		uc := NewUpdateExpenseUseCase(repo, cr)

		input := UpdateExpenseInput{
			ID:         1,
			Amount:     intPtr(200),
			CategoryID: intPtr(2),
			Memo:       "updated memo",
			SpentAt:    "2025-02-01",
			Status:     "confirmed",
		}
		out, err := uc.Execute("test-user", input)

		assert.NoError(t, err)
		assert.True(t, repo.called)
		assert.Equal(t, "confirmed", out.Status)
		assert.Equal(t, 200, out.Amount)
		assert.Equal(t, 2, out.Category.ID)
		assert.Equal(t, "updated memo", out.Memo)
		assert.Equal(t, "2025-02-01", out.SpentAt)
	})

	t.Run("confirmed no status change", func(t *testing.T) {
		t.Parallel()

		repo := &mockUpdateExpenseRepo{current: domain.Expense{ID: 2, Amount: 300, Memo: "c-old", SpentAt: "2025-03-01", Status: "confirmed", Category: domain.Category{ID: 3}}}
		cr := &mockCategoryRepo{exists: map[int32]bool{4: true}}
		uc := NewUpdateExpenseUseCase(repo, cr)

		input := UpdateExpenseInput{
			ID:         2,
			Amount:     intPtr(350),
			CategoryID: intPtr(4),
			Memo:       "c-updated",
			SpentAt:    "2025-03-15",
			Status:     "",
		}
		out, err := uc.Execute("test-user", input)

		assert.NoError(t, err)
		assert.True(t, repo.called)
		assert.Equal(t, "confirmed", out.Status)
		assert.Equal(t, 350, out.Amount)
	})

	t.Run("planned no status change", func(t *testing.T) {
		t.Parallel()

		repo := &mockUpdateExpenseRepo{current: domain.Expense{ID: 3, Amount: 500, Memo: "p-old", SpentAt: "2025-04-01", Status: "planned", Category: domain.Category{ID: 5}}}
		cr := &mockCategoryRepo{exists: map[int32]bool{6: true}}
		uc := NewUpdateExpenseUseCase(repo, cr)

		input := UpdateExpenseInput{
			ID:         3,
			Amount:     intPtr(550),
			CategoryID: intPtr(6),
			Memo:       "p-updated",
			SpentAt:    "2025-04-10",
			Status:     "",
		}
		out, err := uc.Execute("test-user", input)

		assert.NoError(t, err)
		assert.True(t, repo.called)
		assert.Equal(t, "planned", out.Status)
		assert.Equal(t, 550, out.Amount)
	})
}

func TestUpdateExpense_InvalidTransition_ConfirmedToPlanned(t *testing.T) {
	t.Parallel()

	repo := &mockUpdateExpenseRepo{current: domain.Expense{ID: 100, Amount: 1000, Memo: "confirmed item", SpentAt: "2025-05-01", Status: "confirmed", Category: domain.Category{ID: 10}}}
	cr := &mockCategoryRepo{exists: map[int32]bool{11: true}}
	uc := NewUpdateExpenseUseCase(repo, cr)

	input := UpdateExpenseInput{
		ID:         100,
		Amount:     intPtr(1200),
		CategoryID: intPtr(11),
		Memo:       "try revert to planned",
		SpentAt:    "2025-05-02",
		Status:     "planned",
	}

	_, err := uc.Execute("test-user", input)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)
	assert.False(t, repo.called)
}

func TestUpdateExpense_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockUpdateExpenseRepo{getErr: sql.ErrNoRows}
	cr := &mockCategoryRepo{}
	uc := NewUpdateExpenseUseCase(repo, cr)

	input := UpdateExpenseInput{
		ID:         9999,
		Amount:     intPtr(500),
		CategoryID: intPtr(1),
		Memo:       "not found case",
		SpentAt:    "2025-06-01",
		Status:     "planned",
	}

	_, err := uc.Execute("test-user", input)
	assert.ErrorIs(t, err, ErrInvalidStatusTransition)
	assert.False(t, repo.called)
}

// --- mocks for list expenses ---

type mockListExpenseRepo struct {
	fn func(userID string, year, month int) ([]domain.Expense, error)
}

func (m *mockListExpenseRepo) FindByMonth(ctx context.Context, userID string, year, month int) ([]domain.Expense, error) {
	if m.fn != nil {
		return m.fn(userID, year, month)
	}
	return nil, nil
}

func TestListExpenseUseCase_DefaultsToCurrentMonth(t *testing.T) {
	t.Parallel()

	var gotYear, gotMonth int
	repo := &mockListExpenseRepo{
		fn: func(userID string, year, month int) ([]domain.Expense, error) {
			gotYear = year
			gotMonth = month
			return nil, nil
		},
	}
	uc := NewListExpensesUseCase(repo)
	fixed := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	uc.nowFn = func() time.Time { return fixed }

	_, err := uc.Execute(context.Background(), "user1", MonthFilter{})
	require.NoError(t, err)
	assert.Equal(t, 2024, gotYear)
	assert.Equal(t, 6, gotMonth)
}

func TestListExpenseUseCase_WithMonthFilter(t *testing.T) {
	t.Parallel()

	var gotYear, gotMonth int
	repo := &mockListExpenseRepo{
		fn: func(userID string, year, month int) ([]domain.Expense, error) {
			gotYear = year
			gotMonth = month
			return []domain.Expense{{ID: 1}}, nil
		},
	}
	uc := NewListExpensesUseCase(repo)

	result, err := uc.Execute(context.Background(), "user1", MonthFilter{Year: 2024, Month: 11})
	require.NoError(t, err)
	assert.Equal(t, 2024, gotYear)
	assert.Equal(t, 11, gotMonth)
	assert.Len(t, result, 1)
}

func TestListExpenseUseCase_PartialZeroDoesNotOverride(t *testing.T) {
	t.Parallel()

	var gotYear, gotMonth int
	repo := &mockListExpenseRepo{
		fn: func(userID string, year, month int) ([]domain.Expense, error) {
			gotYear = year
			gotMonth = month
			return nil, nil
		},
	}
	uc := NewListExpensesUseCase(repo)

	// year が指定されているが month=0 の場合、year を上書きしない
	_, err := uc.Execute(context.Background(), "user1", MonthFilter{Year: 2024, Month: 0})
	require.NoError(t, err)
	assert.Equal(t, 2024, gotYear)
	assert.Equal(t, 0, gotMonth)
}

func TestListExpenseUseCase_RepoError(t *testing.T) {
	t.Parallel()

	repo := &mockListExpenseRepo{
		fn: func(userID string, year, month int) ([]domain.Expense, error) {
			return nil, errors.New("db error")
		},
	}
	uc := NewListExpensesUseCase(repo)

	_, err := uc.Execute(context.Background(), "user1", MonthFilter{Year: 2024, Month: 1})
	assert.Error(t, err)
}
