package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"money-buddy-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// "Bearer <token>" から token 部分を取得
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		idToken := parts[1]

		// ID Token検証
		token, err := auth.FirebaseAuth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			// エラーログから実際のトークンを除外（セキュリティ対策）
			log.Printf("Failed to verify ID token: %v (token omitted for security)", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに保存
		if token.UID == "" {
			log.Println("Token verification succeeded but UID is empty")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}
		c.Set(string(UserIDKey), token.UID)
		c.Next()
	}
}

// コンテキストからユーザーIDを取得するヘルパー関数
// ミドルウェアを通過している場合は必ずユーザーIDが存在する
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		// このケースは通常発生しないが、念のためログ出力
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		c.Abort()
		return ""
	}
	return userID.(string)
}
