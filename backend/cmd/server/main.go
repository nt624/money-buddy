package main

import (
	"log"
	"os"
	"strings"

	dbgen "pace-wallet-backend/db/generated"
	"pace-wallet-backend/infra/auth"
	"pace-wallet-backend/infra/repository"
	"pace-wallet-backend/infra/transaction"
	"pace-wallet-backend/internal/db"
	"pace-wallet-backend/internal/handlers"
	"pace-wallet-backend/internal/middleware"
	"pace-wallet-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	// Firebase Admin初期化
	if err := auth.InitFirebase(); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// 環境変数から設定を読み込み
	allowedOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:3000")
	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")

	// 本番環境ではリリースモードに設定
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS設定（複数オリジン対応）
	origins := strings.Split(allowedOrigins, ",")
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 許可されたオリジンリストをチェック
		allowed := false
		for _, o := range origins {
			if strings.TrimSpace(o) == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24時間キャッシュ
			c.Writer.Header().Set("Vary", "Origin")                  // 共有キャッシュ対策
		}

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
	expenseRepo := repository.NewExpenseRepositorySQLC(queries)
	categoryRepo := repository.NewCategoryRepositorySQLC(queries)
	userRepo := repository.NewUserRepositorySQLC(queries)
	fixedCostRepo := repository.NewFixedCostRepositorySQLC(queries)
	dashboardRepo := repository.NewDashboardRepositorySQLC(queries)
	txManager := transaction.NewTxManager(dbConn)

	// ユースケース初期化
	createExpenseUC := usecase.NewCreateExpenseUseCase(expenseRepo, categoryRepo)
	listExpensesUC := usecase.NewListExpensesUseCase(expenseRepo)
	deleteExpenseUC := usecase.NewDeleteExpenseUseCase(expenseRepo)
	updateExpenseUC := usecase.NewUpdateExpenseUseCase(expenseRepo, categoryRepo)
	listCategoriesUC := usecase.NewListCategoriesUseCase(categoryRepo)
	initialSetupUC := usecase.NewCompleteInitialSetupUseCase(userRepo, fixedCostRepo, txManager)
	getUserUC := usecase.NewGetUserUseCase(userRepo)
	updateUserUC := usecase.NewUpdateUserSettingsUseCase(userRepo)
	createFixedCostUC := usecase.NewCreateFixedCostUseCase(fixedCostRepo)
	listFixedCostsUC := usecase.NewListFixedCostsUseCase(fixedCostRepo)
	updateFixedCostUC := usecase.NewUpdateFixedCostUseCase(fixedCostRepo)
	deleteFixedCostUC := usecase.NewDeleteFixedCostUseCase(fixedCostRepo)
	getDashboardUC := usecase.NewGetDashboardUseCase(dashboardRepo)

	// 認証不要なエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 認証が必要なエンドポイント（ミドルウェア適用）
	api := r.Group("/")
	api.Use(middleware.AuthMiddleware())
	handlers.Register(api,
		createExpenseUC, listExpensesUC, deleteExpenseUC, updateExpenseUC,
		listCategoriesUC, initialSetupUC, getUserUC, updateUserUC,
		createFixedCostUC, listFixedCostsUC, updateFixedCostUC, deleteFixedCostUC,
		getDashboardUC,
	)

	log.Printf("Server starting on port %s (env: %s)", port, env)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
