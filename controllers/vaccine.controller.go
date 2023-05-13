package controllers

import (
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

func InsertNewVaccineRecord(c *fiber.Ctx) error {
	var request db.Vacunacione

	// Parse request body from JSON to struct
	c.BodyParser(&request)
	// Trim input fields from request body
	purgeInputData(&request)

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}
}
