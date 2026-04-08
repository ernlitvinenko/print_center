package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Debug().Msg(".env file not found, using system environment variables")
	}

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

// getEnv получает переменную окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
