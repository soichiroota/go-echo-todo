package routes

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(server *echo.Echo) {
	server.GET("/", get)
	server.POST("/", post)
	server.GET("/api/", getRequest)
}
