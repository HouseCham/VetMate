package controllers

import (
	"errors"

	"github.com/HouseCham/VetMate/auth"
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

// InsertNewUser inserts a new user into the database
func InsertNewUser(c *fiber.Ctx) error {
	var request db.Usuario
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	purgeInputData(&request)

	// Check if email is already in use
	if message, status, err :=  checkEmailAlreadyInUse(request.Email, true, c); err != nil {
		c.Status(status)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "There is an error with the request",
			"error":   err.Error(),
		})
	}

	// Hash password
	// if error occurs, return 500
	request.Password, err = util.HashPassword(request.Password)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error hashing password",
			"error":   err.Error(),
		})
	}

	params := db.InsertNewUserParams{
		Nombre:       request.Nombre,
		ApellidoP:    request.ApellidoP,
		ApellidoM:    request.ApellidoM,
		Email:        request.Email,
		Telefono:     request.Telefono,
		Password: request.Password,
		Calle:        request.Calle,
		Colonia:      request.Colonia,
		Ciudad:       request.Ciudad,
		Estado:       request.Estado,
		Cp:           request.Cp,
		Pais:         request.Pais,
		NumExt:       request.NumExt,
		NumInt:       request.NumInt,
		Referencia:   request.Referencia,
	}

	// Inserting info into the database
	// if error occurs, return 500
	err = Queries.InsertNewUser(c.Context(), params)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Sorry, there was an error",
			"error":   err.Error(),
		})
	}

	return nil
}
// GetUserByEmail gets a user by email
// if user does not exist, return 404
func GetUserById(c *fiber.Ctx) error {
	return errors.New("not implemented yet")
}

func LoginUser(c *fiber.Ctx) error {
	var request db.Usuario
	var err error

	c.BodyParser(&request)
	purgeInputData(&request)

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err = validations.ValidateRequest(&request, 3); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "There is an error with the request",
			"error":   err.Error(),
		})
	}

	user, err := Queries.GetUserByEmail(c.Context(), request.Email)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Sorry, there was an error loading credentials",
			"error":   err.Error(),
		})
	}

	// Compare password
	// if error occurs, return 500
	if err := util.CheckPassword(request.Password, user.Password); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Wrong password",
			"error":   err.Error(),
		})
	}

	// Generating jwt
	// in case of error, returns 500 with error message
	tokenString, err := auth.GenerateJWT(user.ID, false)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Failed to create the token",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success",
		"token":   tokenString,
	})
}