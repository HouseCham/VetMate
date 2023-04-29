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
	app.Put("/api/v1/vet/update", middleware.JwtMiddleware() , controllers.UpdateVet)
	app.Delete("/api/v1/vet/delete", middleware.JwtMiddleware() , controllers.DeleteVet)

	// User Routes
	app.Post("/api/v1/user/insert", controllers.InsertNewUser)
	app.Get("/api/v1/user/:id", middleware.JwtMiddleware() , controllers.GetUserById)
}