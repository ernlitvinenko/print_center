package auth

import (
	"backend/core/models"
	"backend/core/services"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

var authSvc *services.AuthService

// InitAuthService инициализирует сервис авторизации (вызывать один раз при старте)
func InitAuthService(svc *services.AuthService) {
	authSvc = svc
}

func AddHandlers(instance *echo.Echo) error {
	g := instance.Group("/api/v1/auth")
	g.POST("/login", Login)
	return nil
}

func Login(c *echo.Context) error {
	rd := &models.LoginRequest{}

	if err := c.Bind(rd); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to bind request")
		return c.JSON(http.StatusBadRequest, models.LoginResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Валидация входных данных
	if err := validateLoginRequest(rd); err != nil {
		return c.JSON(http.StatusBadRequest, models.LoginResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	// Очистка номера телефона от лишних символов
	phone := cleanPhoneNumber(*rd.Phone)

	// Аутентификация
	profile, err := authSvc.Authenticate(context.Background(), phone, *rd.Password)
	if err != nil {
		log.Logger.Warn().Err(err).Str("phone", phone).Msg("Login failed")
		return c.JSON(http.StatusUnauthorized, models.LoginResponse{
			Success: false,
			Error:   "Invalid phone number or password",
		})
	}

	// TODO: Здесь будет генерация JWT токена
	// Пока возвращаем пустую строку
	token := ""

	log.Logger.Info().Str("phone", phone).Int32("user_id", profile.ID).Msg("User logged in successfully")

	// Успешный ответ
	return c.JSON(http.StatusOK, models.LoginResponse{
		Success: true,
		Token:   token,
		User: models.UserInfo{
			ID:        profile.ID,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Email:     profile.Email,
			Phone:     phone,
		},
	})
}

// validateLoginRequest проверяет корректность данных для входа
func validateLoginRequest(req *models.LoginRequest) error {
	if req.Phone == nil || strings.TrimSpace(*req.Phone) == "" {
		return fmt.Errorf("phone number is required")
	}

	if req.Password == nil || strings.TrimSpace(*req.Password) == "" {
		return fmt.Errorf("password is required")
	}

	// Проверяем, что телефон содержит только цифры и допустимые символы
	cleanedPhone := cleanPhoneNumber(*req.Phone)
	if !isValidPhone(cleanedPhone) {
		return fmt.Errorf("invalid phone number format")
	}

	// Минимальная длина пароля
	if len(*req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	return nil
}

// cleanPhoneNumber очищает номер телефона от лишних символов
func cleanPhoneNumber(phone string) string {
	// Удаляем все нецифровые символы
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(phone, "")
}

// isValidPhone проверяет корректность номера телефона
func isValidPhone(phone string) bool {
	// Телефон должен содержать от 10 до 15 цифр
	re := regexp.MustCompile(`^\d{10,15}$`)
	return re.MatchString(phone)
}
