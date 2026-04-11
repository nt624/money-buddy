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
	createCategory *usecase.CreateCategoryUseCase,
	updateCategory *usecase.UpdateCategoryUseCase,
	deleteCategory *usecase.DeleteCategoryUseCase,
	reorderCategories *usecase.ReorderCategoriesUseCase,
	initialSetup *usecase.CompleteInitialSetupUseCase,
	getUser *usecase.GetUserUseCase,
	updateUser *usecase.UpdateUserSettingsUseCase,
	createFixedCost *usecase.CreateFixedCostUseCase,
	listFixedCosts *usecase.ListFixedCostsUseCase,
	updateFixedCost *usecase.UpdateFixedCostUseCase,
	deleteFixedCost *usecase.DeleteFixedCostUseCase,
	getDashboard *usecase.GetDashboardUseCase,
	getMonthlySettings *usecase.GetMonthlySettingsUseCase,
	upsertMonthlySettings *usecase.UpsertMonthlySettingsUseCase,
	deleteMonthlySettings *usecase.DeleteMonthlySettingsUseCase,
) {
	r.POST("/expenses", NewCreateExpenseHandler(createExpense).Handle)
	r.GET("/expenses", NewListExpensesHandler(listExpenses).Handle)
	r.DELETE("/expenses/:id", NewDeleteExpenseHandler(deleteExpense).Handle)
	r.PUT("/expenses/:id", NewUpdateExpenseHandler(updateExpense).Handle)

	// categories - /categories/reorder を /:id より先に登録
	r.GET("/categories", NewListCategoriesHandler(listCategories).Handle)
	r.POST("/categories", NewCreateCategoryHandler(createCategory).Handle)
	r.PUT("/categories/reorder", NewReorderCategoriesHandler(reorderCategories).Handle)
	r.PUT("/categories/:id", NewUpdateCategoryHandler(updateCategory).Handle)
	r.DELETE("/categories/:id", NewDeleteCategoryHandler(deleteCategory).Handle)

	r.POST("/setup", NewCompleteInitialSetupHandler(initialSetup).Handle)
	r.GET("/user/me", NewGetUserHandler(getUser).Handle)
	r.PUT("/user/me", NewUpdateUserSettingsHandler(updateUser).Handle)
	r.POST("/fixed-costs", NewCreateFixedCostHandler(createFixedCost).Handle)
	r.GET("/fixed-costs", NewListFixedCostsHandler(listFixedCosts).Handle)
	r.PUT("/fixed-costs/:id", NewUpdateFixedCostHandler(updateFixedCost).Handle)
	r.DELETE("/fixed-costs/:id", NewDeleteFixedCostHandler(deleteFixedCost).Handle)
	r.GET("/dashboard", NewGetDashboardHandler(getDashboard).Handle)

	r.GET("/monthly-settings", NewGetMonthlySettingsHandler(getMonthlySettings).Handle)
	r.PUT("/monthly-settings", NewUpsertMonthlySettingsHandler(upsertMonthlySettings).Handle)
	r.DELETE("/monthly-settings", NewDeleteMonthlySettingsHandler(deleteMonthlySettings).Handle)
}
