package handlers

import (
	"backend/docs"
	"net/http"

	"github.com/labstack/echo/v5"
)

// RegisterSwagger adds a /swagger/* route that serves the swagger UI.
func RegisterSwagger(e *echo.Echo) {
	e.GET("/swagger/*", func(c *echo.Context) error {
		path := c.Param("*")
		switch path {
		case "swagger.json":
			c.Response().Header().Set("Content-Type", "application/json")
			_, err := c.Response().Write([]byte(docs.SwaggerInfo.ReadDoc()))
			return err
		case "":
			fallthrough
		default:
			return c.HTML(http.StatusOK, swaggerHTML)
		}
	})
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Print Center API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      SwaggerUIBundle({
        url: "/swagger/swagger.json",
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout",
        deepLinking: true
      });
    };
  </script>
</body>
</html>`
