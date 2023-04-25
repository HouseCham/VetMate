package controllers

import (
	"errors"

	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
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
	if message, err := checkEmailAlreadyInUse(request.Email, true, c); err != nil {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if isUserValid, err := validations.ValidateUser(request); !isUserValid {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "There is an error with the request",
			"error": err.Error(),
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

func GetUserById(c *fiber.Ctx) error {
	return errors.New("not implemented yet")
}