package main

import (
	"backend/config"
	"backend/core/api/v1/auth"
	"backend/core/repositories"
	"backend/core/services"
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Меняем стандартный логгер на zerolog
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключаемся к базе данных
	conn, err := config.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer conn.Close(context.Background())

	// Инициализируем репозитории
	queries := repositories.New(conn)

	// Инициализируем сервисы
	authService := services.NewAuthService(queries)

	// Инициализируем сервис авторизации
	auth.InitAuthService(authService)

	// Создаем Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Health check
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"message": "Print Center API is running",
		})
	})

	// Регистрируем роуты
	if err := auth.AddHandlers(e); err != nil {
		log.Fatal().Err(err).Msg("Failed to add auth handlers")
	}

	// Запускаем сервер
	log.Info().Str("port", cfg.Port).Msg("Starting server...")
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
