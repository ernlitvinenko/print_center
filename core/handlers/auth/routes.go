package auth

import (
	"backend/core/middleware"
	"backend/core/storage"

	"github.com/labstack/echo/v5"
)

var authStorage storage.AuthStorage
var jwtStorage storage.JWTStorage

// InitStorage связывает handlers с реализациями storage.
func InitStorage(auth storage.AuthStorage, jwt storage.JWTStorage) {
	authStorage = auth
	jwtStorage = jwt
}

// InitializeRouter регистрирует маршруты /auth/*.
func InitializeRouter(app *echo.Echo) {
	g := app.Group("/api/v1/auth")

	g.POST("/login", loginHandler)

	protected := g.Group("")
	protected.GET("/me", middleware.JWT(jwtStorage)(getMeHandler))
}
