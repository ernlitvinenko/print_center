package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// Config содержит конфигурацию приложения
type Config struct {
	DB       *DBConfig
	Port     string
	LogLevel string
}

// DBConfig содержит параметры подключения к БД
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		DB: &DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "print_center"),
		},
		Port:     getEnv("APP_PORT", "8000"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// GetDSN возвращает строку подключения к БД
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}

// ConnectDB устанавливает подключение к базе данных
func ConnectDB(cfg *DBConfig) (*pgx.Conn, error) {
	dsn := cfg.GetDSN()

	log.Info().
		Str("host", cfg.Host).
		Str("port", cfg.Port).
		Str("database", cfg.DBName).
		Msg("Connecting to database...")

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем соединение
	if err := conn.Ping(context.Background()); err != nil {
		conn.Close(context.Background())
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("Successfully connected to database")
	return conn, nil
}

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
