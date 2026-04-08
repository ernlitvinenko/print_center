package enviroment

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Environment хранит все переменные окружения приложения.
type Environment struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	APPHost    string
	APPPort    string
}

var (
	envInstance *Environment
	envOnce     sync.Once
)

// GetInstance загружает .env и возвращает singleton Environment.
func GetInstance() *Environment {
	envOnce.Do(func() {
		if os.Getenv("PROD") != "1" {
			if err := godotenv.Load(".env"); err != nil {
				log.Debug().Msg(".env file not found, using system environment variables")
			}
		}

		envInstance = &Environment{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "postgres"),
			DBPassword: getEnv("DB_PASSWORD", "postgres"),
			DBName:     getEnv("DB_NAME", "print_center"),
			JWTSecret:  getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			APPHost:    getEnv("APP_HOST", "0.0.0.0"),
			APPPort:    getEnv("APP_PORT", "8000"),
		}
	})
	return envInstance
}

// DSN формирует строку подключения к PostgreSQL.
func (e *Environment) DSN() string {
	return "postgres://" + e.DBUser + ":" + e.DBPassword +
		"@" + e.DBHost + ":" + e.DBPort + "/" + e.DBName +
		"?sslmode=disable"
}

// getEnv возвращает переменную окружения или значение по умолчанию.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
