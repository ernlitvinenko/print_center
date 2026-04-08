package main

import (
	"backend/config"
	"backend/core/api/v1/auth"
	"backend/core/repositories"
	"backend/core/services"
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Инициализация логгера
	initLogger()

	// Загрузка конфигурации
	cfg := config.Load()

	// Подключение к БД
	conn := initDatabase(cfg)
	defer conn.Close(context.Background())

	// Инициализация сервисов
	queries := repositories.New(conn)
	jwtConfig := services.NewJWTConfig()
	authService := services.NewAuthService(queries, jwtConfig)
	auth.InitAuthService(authService, jwtConfig)

	// Настройка роутера
	e := initEcho()

	// Регистрация роутов
	registerRoutes(e)

	// Запуск сервера
	log.Info().Str("port", cfg.Port).Msg("Starting server...")
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func initLogger() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func initDatabase(cfg *config.Config) *pgx.Conn {
	conn, err := config.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	return conn
}

func initEcho() *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Health check
	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"message": "Print Center API is running",
		})
	})

	return e
}

func registerRoutes(e *echo.Echo) {
	if err := auth.AddHandlers(e); err != nil {
		log.Fatal().Err(err).Msg("Failed to add auth handlers")
	}
}
