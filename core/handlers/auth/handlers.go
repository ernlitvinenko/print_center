package auth

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

// LoginRequest — запрос на авторизацию.
type LoginRequest struct {
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
}

// LoginResponse — ответ после авторизации.
type LoginResponse struct {
	Success bool      `json:"success"`
	Token   string    `json:"token,omitempty"`
	User    *UserInfo `json:"user,omitempty"`
	Error   string    `json:"error,omitempty"`
}

// UserInfo — данные пользователя.
type UserInfo struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// loginHandler обрабатывает POST /login.
// @Summary Авторизация пользователя
// @Description Вход по номеру телефона и паролю, возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Данные для входа"
// @Success 200 {object} LoginResponse "Успешный вход"
// @Failure 400 {object} LoginResponse "Ошибка валидации"
// @Failure 401 {object} LoginResponse "Неверные учётные данные"
// @Router /auth/login [post]
func loginHandler(c *echo.Context) error {
	req := &LoginRequest{}
	if err := c.Bind(req); err != nil {
		log.Logger.Error().Err(err).Msg("Failed to bind login request")
		return c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if err := validateLoginRequest(req); err != nil {
		return c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	phone := cleanPhoneNumber(*req.Phone)
	profile, token, err := authStorage.Authenticate(c.Request().Context(), phone, *req.Password)
	if err != nil {
		log.Logger.Warn().Err(err).Str("phone", phone).Msg("Login failed")
		return c.JSON(http.StatusUnauthorized, LoginResponse{
			Success: false,
			Error:   "Invalid phone number or password",
		})
	}

	log.Logger.Info().Str("phone", phone).Int32("user_id", profile.ID).Msg("User logged in successfully")

	return c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Token:   token,
		User: &UserInfo{
			ID:        profile.ID,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Email:     profile.Email,
			Phone:     phone,
		},
	})
}

// getMeHandler возвращает данные текущего пользователя.
// @Summary Данные текущего пользователя
// @Description Возвращает информацию о пользователе по JWT токену
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} LoginResponse "Данные пользователя"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Router /auth/me [get]
func getMeHandler(c *echo.Context) error {
	userID := c.Get("user_id").(int32)
	phone := c.Get("phone").(string)
	firstName := c.Get("first_name").(string)
	lastName := c.Get("last_name").(string)

	return c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		User: &UserInfo{
			ID:        userID,
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phone,
		},
	})
}

// validateLoginRequest проверяет обязательные поля.
func validateLoginRequest(req *LoginRequest) error {
	if req.Phone == nil || strings.TrimSpace(*req.Phone) == "" {
		return &ValidationError{"phone number is required"}
	}
	if req.Password == nil || strings.TrimSpace(*req.Password) == "" {
		return &ValidationError{"password is required"}
	}

	cleanedPhone := cleanPhoneNumber(*req.Phone)
	if !isValidPhone(cleanedPhone) {
		return &ValidationError{"invalid phone number format"}
	}
	if len(*req.Password) < 6 {
		return &ValidationError{"password must be at least 6 characters"}
	}
	return nil
}

// cleanPhoneNumber удаляет все нецифровые символы.
func cleanPhoneNumber(phone string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(phone, "")
}

// isValidPhone проверяет формат телефона (10-15 цифр).
func isValidPhone(phone string) bool {
	re := regexp.MustCompile(`^\d{10,15}$`)
	return re.MatchString(phone)
}

// ValidationError — ошибка валидации входных данных.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
