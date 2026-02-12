package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware_NoAuthorizationHeader(t *testing.T) {
	router := setupTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "認証ヘッダーが必要です")
}

func TestAuthMiddleware_InvalidAuthorizationFormat_NoBearer(t *testing.T) {
	router := setupTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidToken")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "認証形式が正しくありません")
}

func TestAuthMiddleware_InvalidAuthorizationFormat_OnlyBearer(t *testing.T) {
	router := setupTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "認証形式が正しくありません")
}

func TestAuthMiddleware_InvalidAuthorizationFormat_BearerWithSpacesOnly(t *testing.T) {
	router := setupTestRouter()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer   ")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "認証形式が正しくありません")
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// 注意: このテストは実際のFirebase Authが初期化されている必要があります
	// Firebase Authが未初期化の場合、トークン検証でpanicが発生します
	// 統合テストとして別途実装することを推奨します
	t.Skip("このテストは実際のFirebase Authが必要なため、スキップします")
}

func TestAuthMiddleware_EmptyUID(t *testing.T) {
	// 注意: このテストは実際のFirebase Authを使用する必要があるため、
	// 統合テストとして別途実装することを推奨します。
	// Firebase Authは通常空のUIDを返さないため、このケースはモックでのみテスト可能です。
	t.Skip("このテストはFirebase Authのモックが必要なため、スキップします")
}

func TestAuthMiddleware_Success(t *testing.T) {
	// 注意: このテストは実際の有効なFirebaseトークンが必要です。
	// CI/CD環境では、テスト用のFirebaseプロジェクトと有効なトークンを使用してください。
	t.Skip("このテストは有効なFirebaseトークンが必要なため、スキップします")

	// 実装例（有効なトークンがある場合）:
	// router := setupTestRouter()
	// router.Use(AuthMiddleware())
	// router.GET("/test", func(c *gin.Context) {
	//     userID, _ := GetUserID(c)
	//     c.JSON(http.StatusOK, gin.H{"user_id": userID})
	// })
	//
	// w := httptest.NewRecorder()
	// req, _ := http.NewRequest("GET", "/test", nil)
	// req.Header.Set("Authorization", "Bearer VALID_FIREBASE_TOKEN")
	// router.ServeHTTP(w, req)
	//
	// assert.Equal(t, http.StatusOK, w.Code)
	// assert.Contains(t, w.Body.String(), "user_id")
}

func TestGetUserID_UserIDExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUserID := "test-user-123"
	c.Set(string(UserIDKey), expectedUserID)

	userID, ok := GetUserID(c)

	assert.True(t, ok)
	assert.Equal(t, expectedUserID, userID)
}

func TestGetUserID_UserIDNotExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// UserIDをセットしない

	userID, ok := GetUserID(c)

	assert.False(t, ok)
	assert.Equal(t, "", userID)
}

func TestGetUserID_AfterMiddleware(t *testing.T) {
	// ミドルウェアを通過した後のGetUserIDの動作確認
	router := setupTestRouter()

	// ミドルウェアの代わりにUserIDを直接セット
	router.Use(func(c *gin.Context) {
		c.Set(string(UserIDKey), "middleware-user-123")
		c.Next()
	})

	router.GET("/test", func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの取得に失敗しました"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "middleware-user-123")
}
