package helpers

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// GetEnv возвращает значение переменной окружения или defaultValue.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt парсит переменную окружения в int. При ошибке возвращает defaultValue.
func GetEnvInt(key string, defaultValue int) int {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}

	var val int
	if _, err := sscanf(str, "%d", &val); err != nil {
		log.Warn().Str("key", key).Str("value", str).Msg("Failed to parse int from env, using default")
		return defaultValue
	}
	return val
}

// sscanf — минимальная замена fmt.Sscanf для одного аргумента.
func sscanf(s, format string, a ...interface{}) (int, error) {
	return fmt.Sscanf(s, format, a...)
}
