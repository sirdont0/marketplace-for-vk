package auth_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"marketplace/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	// Устанавливаем тестовые переменные окружения
	os.Setenv("JWT_SECRET", "test_secret_1234567890")
	os.Setenv("JWT_EXPIRE_HOURS", "1")

	t.Run("Generate and parse valid token", func(t *testing.T) {
		userID := 123
		token, err := auth.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := auth.ParseToken(token)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%d", userID), claims.Subject)
	})

	t.Run("Invalid token", func(t *testing.T) {
		_, err := auth.ParseToken("invalid.token.here")
		assert.Error(t, err)
	})

	t.Run("Expired token", func(t *testing.T) {
		// Создаем просроченный токен
		claims := jwt.RegisteredClaims{
			Subject:   "123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte("test_secret_1234567890"))

		_, err := auth.ParseToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})
}
