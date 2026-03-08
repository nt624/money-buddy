package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"money-buddy-backend/internal/domain"
	"money-buddy-backend/internal/middleware"
)

type listExpensesUseCase interface {
	Execute(userID string) ([]domain.Expense, error)
}

type ListExpensesHandler struct {
	uc listExpensesUseCase
}

func NewListExpensesHandler(uc listExpensesUseCase) *ListExpensesHandler {
	return &ListExpensesHandler{uc: uc}
}

func (h *ListExpensesHandler) Handle(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
		return
	}

	expenses, err := h.uc.Execute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支出の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}
