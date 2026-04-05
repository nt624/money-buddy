package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type updateExpenseUseCase interface {
	Execute(userID string, input usecase.UpdateExpenseInput) (domain.Expense, error)
}

type UpdateExpenseHandler struct {
	uc updateExpenseUseCase
}

func NewUpdateExpenseHandler(uc updateExpenseUseCase) *UpdateExpenseHandler {
	return &UpdateExpenseHandler{uc: uc}
}

func (h *UpdateExpenseHandler) Handle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "支出IDが正しくありません"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "支出IDが正しくありません"})
		return
	}

	type updateBody struct {
		Amount     *int   `json:"amount"`
		CategoryID *int   `json:"category_id"`
		Memo       string `json:"memo"`
		SpentAt    string `json:"spent_at"`
		Status     string `json:"status"`
	}
	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateExpenseInput{
		ID:         int(id),
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

	exp, err := h.uc.Execute(userID, input)
	if err != nil {
		var ve *usecase.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ve.Message})
			return
		}
		if errors.Is(err, usecase.ErrInvalidStatusTransition) {
			c.JSON(http.StatusConflict, gin.H{"error": "ステータスの変更ができません（確定済みは予定に戻せません）"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバーエラーが発生しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expense": exp})
}
