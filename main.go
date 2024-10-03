package main

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/create-go-app/fiber-go-template/app/dao"
	"os"

	"github.com/create-go-app/fiber-go-template/pkg/configs"
	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/create-go-app/fiber-go-template/pkg/routes"
	"github.com/create-go-app/fiber-go-template/pkg/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/create-go-app/fiber-go-template/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// dao
	mainDao, err := dao.NewDao()
	if err != nil {
		panic(err)
	}

	// services
	service := controllers.NewService(mainDao)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app)           // Register a route for API Docs (Swagger).
	routes.PublicRoutes(*service, app) // Register a public routes for app.
	routes.NotFoundRoute(app)          // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
