package main

import (
	"backend/core/handlers"
	"backend/core/helpers"
	"backend/core/tasks"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	app := echo.New()

	// Logger
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// CORS
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type", "Accept", "Origin"},
	}))

	app.Use(middleware.RequestLogger())
	app.Use(middleware.Recover())

	// Health check
	app.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"message": "Print Center API is running",
		})
	})

	// Startup tasks (init env, db, storage)
	tasks.OnStartup()
	defer tasks.OnShutdown()

	// Init storage and handlers
	handlers.InitStorage()

	// Register routes
	handlers.InitializeRoutes(app)

	// Start server
	port := helpers.GetEnv("APP_PORT", "8000")
	log.Info().Str("port", port).Msg("Starting server...")

	go func() {
		if err := app.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	helpers.GracefulShutdown(app)
}
