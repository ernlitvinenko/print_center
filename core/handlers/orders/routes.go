package orders

import (
	"backend/core/middleware"

	"github.com/labstack/echo/v5"
)

// InitializeRouter регистрирует маршруты /orders/*.
func InitializeRouter(app *echo.Echo) {
	g := app.Group("/api/v1/orders", middleware.JWT(jwtStorage))

	g.GET("", getOrdersHandler)
	g.POST("", createOrderHandler)
}
