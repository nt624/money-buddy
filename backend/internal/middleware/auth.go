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
		c.Set(string(UserIDKey), token.UID)
		c.Next()
	}
}

// コンテキストからユーザーIDを取得するヘルパー関数
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get(string(UserIDKey))
	if !exists {
		return ""
	}
	return userID.(string)
}
