package services

import (
	"backend/core/repositories"
	"context"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repositories.Queries
	jwtConfig *JWTConfig
}

func NewAuthService(repo *repositories.Queries, jwtConfig *JWTConfig) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

// HashPassword генерирует bcrypt хеш из пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword сравнивает пароль с хешем
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Authenticate проверяет учетные данные пользователя и возвращает профиль + токен
func (s *AuthService) Authenticate(ctx context.Context, phone string, password string) (*repositories.Profile, string, error) {
	// Преобразуем телефон в числовой формат
	phoneDgt, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		return nil, "", fmt.Errorf("invalid phone number format: %w", err)
	}

	// Получаем профиль из базы данных
	profile, err := s.repo.GetProfile(ctx, repositories.GetProfileParams{
		Email:    "",
		PhoneDgt: phoneDgt,
	})
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	// Проверяем пароль
	if err := CheckPassword(password, profile.Password); err != nil {
		return nil, "", fmt.Errorf("invalid password: %w", err)
	}

	// Генерируем JWT токен
	token, err := s.jwtConfig.GenerateToken(profile.ID, phone, profile.FirstName, profile.LastName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return &profile, token, nil
}
