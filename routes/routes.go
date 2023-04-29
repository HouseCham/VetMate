package routes

import (
	"github.com/HouseCham/VetMate/controllers"
	"github.com/HouseCham/VetMate/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetAllRoutes is a function that sets up all routes for the application
func SetAllRoutes(app *fiber.App) {
	// Vet Routes
	app.Get("/api/v1/vet/get", middleware.JwtMiddleware() , controllers.GetVetById)
	app.Post("/api/v1/vet", controllers.InsertNewVet)
	app.Post("/api/v1/vet/login", controllers.LoginVet)
	app.Put("/api/v1/vet/update", middleware.JwtMiddleware() , controllers.UpdateVet)
	app.Delete("/api/v1/vet/delete", middleware.JwtMiddleware() , controllers.DeleteVet)

	// User Routes
	app.Get("/api/v1/user/:id", middleware.JwtMiddleware() , controllers.GetUserById)
	app.Post("/api/v1/user/insert", controllers.InsertNewUser)
	app.Post("/api/v1/user/login", controllers.LoginUser)
	app.Put("/api/v1/user/update", middleware.JwtMiddleware() , controllers.UpdateUser)
}