package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/middleware"
	"money-buddy-backend/internal/usecase"
)

type createExpenseUseCase interface {
	Execute(userID string, input usecase.CreateExpenseInput) (domain.Expense, error)
}

type CreateExpenseHandler struct {
	uc createExpenseUseCase
}

func NewCreateExpenseHandler(uc createExpenseUseCase) *CreateExpenseHandler {
	return &CreateExpenseHandler{uc: uc}
}

func (h *CreateExpenseHandler) Handle(c *gin.Context) {
	var input usecase.CreateExpenseInput

	type createBody struct {
		Amount     *int   `json:"amount"`
		CategoryID *int   `json:"category_id"`
		Memo       string `json:"memo"`
		SpentAt    string `json:"spent_at"`
		Status     string `json:"status"`
	}
	var body createBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input = usecase.CreateExpenseInput{
		Amount:     body.Amount,
		CategoryID: body.CategoryID,
		Memo:       body.Memo,
		SpentAt:    body.SpentAt,
		Status:     body.Status,
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	expense, err := h.uc.Execute(userID, input)
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"expense": expense})
}
