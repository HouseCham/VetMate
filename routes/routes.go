package routes

import (
	"github.com/HouseCham/VetMate/controllers"
	"github.com/gofiber/fiber/v2"
)

// SetAllRoutes is a function that sets up all routes for the application
func SetAllRoutes(app *fiber.App) {
	app.Get("/", controllers.HelloWorld)

	// Vet Routes
	app.Post("/api/v1/vet", controllers.InsertNewVet)
}