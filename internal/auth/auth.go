package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GetJWTSecret возвращает секретный ключ из переменных окружения
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set in environment variables")
	}
	return secret
}

// JWTExpireDuration возвращает время жизни токена из переменных окружения
func JWTExpireDuration() time.Duration {
	hoursStr := os.Getenv("JWT_EXPIRE_HOURS")
	if hoursStr == "" {
		return 24 * time.Hour // Значение по умолчанию
	}

	hours, err := time.ParseDuration(hoursStr + "h")
	if err != nil {
		return 24 * time.Hour
	}
	return hours
}

// GenerateToken создает новый JWT токен для пользователя
func GenerateToken(userID int) (string, error) {
	// Создаем claims с данными пользователя
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID), // ID пользователя в строковом формате
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTExpireDuration())),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// Создаем токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(GetJWTSecret()))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseToken проверяет и парсит JWT токен
func ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	// Парсим токен с проверкой подписи
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetJWTSecret()), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Проверяем claims
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
