package controllers

import (
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/gofiber/fiber/v2"
)

func InsertNewUser(c *fiber.Ctx) error {
	var request db.InsertNewUserParams
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	trimInputFields(&request)

	// Check if email is already in use
	//! TODO: generate method to check email in users
	if message, err := checkVetEmailAlreadyInUse(request.Email, c); err != nil {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Hash password
	// if error occurs, return 500
	request.PasswordHash, err = util.HashPassword(request.PasswordHash)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error hashing password",
			"error":   err.Error(),
		})
	}

	//! TODO: create validate function for user
	//! GENERATE interface in order to use the same function of validation for both, user and vet

	// Inserting info into the database
	// if error occurs, return 500
	err = Queries.InsertNewUser(c.Context(), request)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Sorry, there was an error",
			"error": err.Error(),
		})
	}

	return nil
}