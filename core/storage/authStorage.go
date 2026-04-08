package storage

import (
	"backend/core/repositories"
	"context"
	"fmt"
	"strconv"
)

// AuthStorage — хранилище для аутентификации и работы с профилями.
type AuthStorage interface {
	Storage
	Authenticate(ctx context.Context, phone, password string) (*repositories.Profile, string, error)
	GetProfileByPhone(ctx context.Context, phone string) (*repositories.Profile, error)
	ListProfiles(ctx context.Context) ([]repositories.Profile, error)
}

// Authenticate проверяет телефон/пароль и возвращает профиль + JWT-токен.
func (s *storage) Authenticate(ctx context.Context, phone, password string) (*repositories.Profile, string, error) {
	phoneDgt, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		return nil, "", fmt.Errorf("invalid phone number format: %w", err)
	}

	profile, err := s.queries.GetProfile(ctx, repositories.GetProfileParams{
		Email:    "",
		PhoneDgt: phoneDgt,
	})
	if err != nil {
		return nil, "", fmt.Errorf("user not found: %w", err)
	}

	if err := checkPassword(password, profile.Password); err != nil {
		return nil, "", fmt.Errorf("invalid password: %w", err)
	}

	token, err := s.GenerateToken(profile.ID, phone, profile.FirstName, profile.LastName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return &profile, token, nil
}

// GetProfileByPhone ищет профиль по номеру телефона.
func (s *storage) GetProfileByPhone(ctx context.Context, phone string) (*repositories.Profile, error) {
	phoneDgt, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number format: %w", err)
	}

	profile, err := s.queries.GetProfile(ctx, repositories.GetProfileParams{
		Email:    "",
		PhoneDgt: phoneDgt,
	})
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &profile, nil
}

// ListProfiles возвращает все профили.
func (s *storage) ListProfiles(ctx context.Context) ([]repositories.Profile, error) {
	return s.queries.ListProfiles(ctx)
}
