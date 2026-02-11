package main

import (
	"log"

	dbgen "money-buddy-backend/db/generated"
	"money-buddy-backend/infra/repository"
	"money-buddy-backend/internal/auth"
	"money-buddy-backend/internal/db"
	"money-buddy-backend/internal/handlers"
	"money-buddy-backend/internal/middleware"
	"money-buddy-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Firebase Admin初期化
	if err := auth.InitFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	dbConn, err := db.NewDB()
	if err != nil {
		panic(err)
	}

	queries := dbgen.New(dbConn)
	repo := repository.NewExpenseRepositorySQLC(queries)
	categoryRepo := repository.NewCategoryRepositorySQLC(queries)
	userRepo := repository.NewUserRepositorySQLC(queries)
	fixedCostRepo := repository.NewFixedCostRepositorySQLC(queries)
	txManager := db.NewSQLTxManager(dbConn)
	dashboardRepo := repository.NewDashboardRepositorySQLC(queries)

	// サービス初期化
	service := services.NewExpenseService(repo, categoryRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	initialSetupService := services.NewInitialSetupService(userRepo, fixedCostRepo, txManager)
	userService := services.NewUserService(userRepo)
	fixedCostService := services.NewFixedCostService(fixedCostRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	// 認証が必要なエンドポイント（ミドルウェア適用）
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	{
		handlers.NewExpenseHandler(api, service)
		handlers.NewCategoryHandler(api, categoryService)
		handlers.NewInitialSetupHandler(api, initialSetupService)
		handlers.NewUserHandler(api, userService)
		handlers.NewFixedCostHandler(api, fixedCostService)
		handlers.NewDashboardHandler(api, dashboardService)
	}

	r.Run() // デフォルトで:8080で起動
}
