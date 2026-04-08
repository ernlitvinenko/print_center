package storage

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// JWTStorage отвечает за генерацию и валидацию JWT-токенов.
type JWTStorage interface {
	Storage
	GenerateToken(userID int32, phone, firstName, lastName string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

// JWTClaims — кастомные claims токена.
type JWTClaims struct {
	UserID    int32  `json:"user_id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	jwt.RegisteredClaims
}

// GenerateToken создаёт подписанный JWT-токен.
func (s *storage) GenerateToken(userID int32, phone, firstName, lastName string) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Phone:     phone,
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "print-center-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.settings.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken проверяет подпись и срок действия токена, возвращает claims.
func (s *storage) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.settings.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		log.Warn().Msg("Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
