package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/domain"
	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"
)

type listExpensesUseCase interface {
	Execute(ctx context.Context, userID string, filter usecase.MonthFilter) ([]domain.Expense, error)
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

	var filter usecase.MonthFilter

	if yearStr := c.Query("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil || year < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "yearは正の整数で指定してください"})
			return
		}
		filter.Year = year
	}

	if monthStr := c.Query("month"); monthStr != "" {
		month, err := strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "monthは1から12の整数で指定してください"})
			return
		}
		filter.Month = month
	}

	if (filter.Year == 0) != (filter.Month == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "yearとmonthは両方指定するか、両方省略してください"})
		return
	}

	expenses, err := h.uc.Execute(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "支出の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}
