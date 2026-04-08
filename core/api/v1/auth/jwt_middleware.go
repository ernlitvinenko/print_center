package auth

import (
	"backend/core/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

// JWTMiddleware middleware для проверки JWT токенов
type JWTMiddleware struct {
	jwtConfig *services.JWTConfig
}

// NewJWTMiddleware создаёт новый middleware для JWT
func NewJWTMiddleware(jwtConfig *services.JWTConfig) *JWTMiddleware {
	return &JWTMiddleware{
		jwtConfig: jwtConfig,
	}
}

// Handle middleware для проверки токена
func (m *JWTMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "missing authorization header",
			})
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid authorization format, expected: Bearer <token>",
			})
		}

		tokenString := parts[1]

		// Валидируем токен
		claims, err := m.jwtConfig.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid or expired token",
			})
		}

		// Сохраняем claims в контекст для использования в обработчиках
		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Set("first_name", claims.FirstName)
		c.Set("last_name", claims.LastName)

		return next(c)
	}
}
