package helpers

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/rs/zerolog/log"
)

// GracefulShutdown ожидает сигнал завершения и останавливает Echo-сервер.
func GracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("Shutting down server...")
}
