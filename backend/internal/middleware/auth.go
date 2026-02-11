package middleware

import (
	"context"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに保存
		if token.UID == "" {
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
