// Package handlers регистрирует все HTTP-маршруты API v1.
package handlers

import (
	"backend/core/handlers/auth"
	"backend/core/handlers/orders"
	"backend/core/storage"

	"github.com/labstack/echo/v5"
)

// InitializeRoutes подключает все маршруты.
func InitializeRoutes(app *echo.Echo) {
	auth.InitializeRouter(app)
	orders.InitializeRouter(app)
}

// InitStorage инициализирует storage для всех handlers.
func InitStorage() {
	authStore := storage.GetInstance[storage.AuthStorage]()
	jwtStore := storage.GetInstance[storage.JWTStorage]()
	orderStore := storage.GetInstance[storage.OrderStorage]()

	auth.InitStorage(authStore, jwtStore)
	orders.InitStorage(orderStore, jwtStore)
}
