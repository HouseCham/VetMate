package routes

import (
	"github.com/HouseCham/VetMate/controllers"
	"github.com/HouseCham/VetMate/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetAllRoutes is a function that sets up all routes for the application
func SetAllRoutes(app *fiber.App) {
	// Vet Routes
	app.Post("/api/v1/vet", controllers.InsertNewVet)
	app.Post("/api/v1/vet/login", controllers.LoginVet)
	app.Get("/api/v1/vet/get", middleware.JwtMiddleware() , controllers.GetVetById)

	// User Routes
	app.Post("/api/v1/user", controllers.InsertNewUser)
	app.Get("/api/v1/user/:id", middleware.JwtMiddleware() , controllers.GetUserById)
}