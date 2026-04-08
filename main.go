package main

import (
	"backend/config"
	v1auth "backend/core/api/v1/auth"
	v1order "backend/core/api/v1/order"
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
	orderService := services.NewOrderService(queries, conn)

	v1auth.InitAuthService(authService, jwtConfig)
	v1order.InitOrderService(orderService)

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

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type", "Accept", "Origin"},
	}))

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
	// Auth routes (публичные + защищённые)
	if err := v1auth.AddHandlers(e); err != nil {
		log.Fatal().Err(err).Msg("Failed to add auth handlers")
	}

	// Order routes (защищённые)
	if err := v1order.AddHandlers(e, v1auth.GetJWTMiddleware()); err != nil {
		log.Fatal().Err(err).Msg("Failed to add order handlers")
	}
}
