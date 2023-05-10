package routes

import (
	"github.com/HouseCham/VetMate/controllers"
	"github.com/HouseCham/VetMate/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetAllRoutes is a function that sets up all routes for the application
func SetAllRoutes(app *fiber.App) {
	// Vet Routes
	app.Post("/api/v1/vet/login", controllers.LoginVet)
	app.Post("/api/v1/vet", controllers.InsertNewVet)
	app.Get("/api/v1/vet/get", middleware.JwtMiddleware(true), controllers.GetVetById)
	app.Put("/api/v1/vet/update", middleware.JwtMiddleware(true), controllers.UpdateVet)
	app.Delete("/api/v1/vet/delete", middleware.JwtMiddleware(true), controllers.DeleteVet)

	// User Routes
	app.Post("/api/v1/user/login", controllers.LoginUser)
	app.Post("/api/v1/user/insert", controllers.InsertNewUser)
	app.Get("/api/v1/user/:id", middleware.JwtMiddleware(false), controllers.GetUserById)
	app.Put("/api/v1/user/update", middleware.JwtMiddleware(false), controllers.UpdateUser)
	app.Delete("/api/v1/user/delete", middleware.JwtMiddleware(false), controllers.DeleteUser)

	// Pet Routes
	app.Post("/api/v1/pet/insertFromUser", middleware.JwtMiddleware(false), controllers.InsertNewPetByUser)
	app.Get("/api/v1/pet/get/:petId", middleware.JwtMiddleware(false), controllers.GetPet)
	app.Put("/api/v1/pet/update/:petId", middleware.JwtMiddleware(false), controllers.UpdatePet)

	// vaccine Routes
	app.Post("/api/v1/vaccine/insert", middleware.JwtMiddleware(true), controllers.InsertNewVaccine)
}