package main

import (
	"log"
	"os"
	"strings"

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
