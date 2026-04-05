package handlers

import (
	"github.com/gin-gonic/gin"

	"pace-wallet-backend/internal/usecase"
)

// Register はすべてのAPIルートを登録します。
func Register(
	r gin.IRouter,
	createExpense *usecase.CreateExpenseUseCase,
	listExpenses *usecase.ListExpensesUseCase,
	deleteExpense *usecase.DeleteExpenseUseCase,
	updateExpense *usecase.UpdateExpenseUseCase,
	listCategories *usecase.ListCategoriesUseCase,
	initialSetup *usecase.CompleteInitialSetupUseCase,
	getUser *usecase.GetUserUseCase,
	updateUser *usecase.UpdateUserSettingsUseCase,
	createFixedCost *usecase.CreateFixedCostUseCase,
	listFixedCosts *usecase.ListFixedCostsUseCase,
	updateFixedCost *usecase.UpdateFixedCostUseCase,
	deleteFixedCost *usecase.DeleteFixedCostUseCase,
	getDashboard *usecase.GetDashboardUseCase,
) {
	r.POST("/expenses", NewCreateExpenseHandler(createExpense).Handle)
	r.GET("/expenses", NewListExpensesHandler(listExpenses).Handle)
	r.DELETE("/expenses/:id", NewDeleteExpenseHandler(deleteExpense).Handle)
	r.PUT("/expenses/:id", NewUpdateExpenseHandler(updateExpense).Handle)
	r.GET("/categories", NewListCategoriesHandler(listCategories).Handle)
	r.POST("/setup", NewCompleteInitialSetupHandler(initialSetup).Handle)
	r.GET("/user/me", NewGetUserHandler(getUser).Handle)
	r.PUT("/user/me", NewUpdateUserSettingsHandler(updateUser).Handle)
	r.POST("/fixed-costs", NewCreateFixedCostHandler(createFixedCost).Handle)
	r.GET("/fixed-costs", NewListFixedCostsHandler(listFixedCosts).Handle)
	r.PUT("/fixed-costs/:id", NewUpdateFixedCostHandler(updateFixedCost).Handle)
	r.DELETE("/fixed-costs/:id", NewDeleteFixedCostHandler(deleteFixedCost).Handle)
	r.GET("/dashboard", NewGetDashboardHandler(getDashboard).Handle)
}
