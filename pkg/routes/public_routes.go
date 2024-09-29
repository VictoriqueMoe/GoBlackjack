package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(s controllers.Service, app *fiber.App) {
	routes := s.GetAllRoutes()
	router := app.Group("/api/v1")
	for _, pRoute := range routes {
		pRoute(router)
	}
}
