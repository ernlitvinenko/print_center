package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig конфигурация для JWT токенов
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// JWTClaims кастомные claims для токена
type JWTClaims struct {
	UserID    int32  `json:"user_id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

// NewJWTConfig создаёт конфигурацию JWT со значениями по умолчанию
func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		TokenDuration: time.Hour * 24, // 24 часа
	}
}

// GenerateToken генерирует JWT токен для пользователя
func (c *JWTConfig) GenerateToken(userID int32, phone, firstName, lastName string) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Phone:     phone,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "print-center-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken проверяет валидность JWT токена и возвращает claims
func (c *JWTConfig) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(c.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
