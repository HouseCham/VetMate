package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/HouseCham/VetMate/auth"
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

// InsertNewVet is a function that inserts a new vet
// to the database
func InsertNewVet(c *fiber.Ctx) error {

	var request db.Veterinario
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	purgeInputData(&request)

	// Check if email is already in use
	// set false as second parameter in order to check for vet emails
	if message, status, err := checkEmailAlreadyInUse(request.Email, false, c); err != nil {
		c.Status(status)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Validate request body
	// if not valid, return 400
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
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

	// Mapping request body to InsertNewVetParams struct
	params := db.InsertNewVetParams{
		Nombre:       request.Nombre,
		ApellidoP:    request.ApellidoP,
		ApellidoM:    request.ApellidoM,
		Password: request.Password,
		Email:        request.Email,
		Telefono:     request.Telefono,
	}

	// Starting transaction
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error starting transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback()

	// getting queries with transaction
	qtx := Queries.WithTx(tx)
	err = qtx.InsertNewVet(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error inserting new vet",
			"error":   err.Error(),
		})
	}
	return tx.Commit()
}

//! TODO: think if we need transaction for selecting data
// GetVetById is a function that gets the vet info
// by the vet id from the fiber context
func GetVetById(c *fiber.Ctx) error {
	// Get the variable from the request context
	// Variable not found or not of type string
	vetId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// then, we need to get the vet info from the database
	// if error occurs, return 404
	mainInfo, err := Queries.GetVetMainInfoById(c.Context(), vetId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not get vet info",
			"error":   err.Error(),
		})
	}
	return c.JSON(mainInfo)
}

// UpdateVet is a function that updates the vet info
// by the vet id from the fiber context
func UpdateVet(c *fiber.Ctx) error {
	// Get the vetId from the request context
	// If variable not found or not of type string
	// return 500 with error message
	vetId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	var request db.Veterinario

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	request.ID = vetId

	// Trim input fields from request body
	purgeInputData(&request)

	// Validate request body
	// if not valid, return 400
	if err := validations.ValidateRequest(&request, 2); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Mapping request body to db.UpdateVetParams struct
	params := db.UpdateVetParams{
		ID:        request.ID,
		Nombre:    request.Nombre,
		ApellidoP: request.ApellidoP,
		ApellidoM: request.ApellidoM,
		Telefono:  request.Telefono,
		ImgUrl:    request.ImgUrl,
	}

	// Starting transaction
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error starting transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback()

	// getting queries with transaction
	qtx := Queries.WithTx(tx)
	err = qtx.UpdateVet(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error updating vet",
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}

// DeleteVet is a function that deletes the vet info
// by the vet id from fiber context
func DeleteVet(c *fiber.Ctx) error {
	// Get the vetId from the request context
	// If variable not found or not of type string
	// return 500 with error message
	vetId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Starting transaction
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error starting transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback()

	// getting queries with transaction
	qtx := Queries.WithTx(tx)
	err = qtx.DeleteVet(c.Context(), vetId)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error deleting vet",
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Implementing Trim() function for LoginRequest struct
func (loginRequest *LoginRequest) Trim() {
	loginRequest.Email = strings.TrimSpace(loginRequest.Email)
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)
}

// LoginVet is a function that logs in the vet
// by checking the email and password
func LoginVet(c *fiber.Ctx) error {
	var request db.Veterinario
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	purgeInputData(&request)

	// Validate request body
	// if not valid, return 400
	if err := validations.ValidateRequest(&request, 3); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	// Get vet info from database
	vet, err := Queries.GetVetByEmail(c.Context(), request.Email)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not get vet info",
			"error":   err.Error(),
		})
	}
	// Compare password
	// if error occurs, return 500
	if err := util.CheckPassword(request.Password, vet.Password); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Wrong password",
			"error":   err.Error(),
		})
	}

	// Generating jwt
	// in case of error, returns 500 with error message
	tokenString, err := auth.GenerateJWT(vet.ID, true)
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

// getIdFromRequestContext is a function that gets the vet id
// from the request context -> c.Locals("userId")
// returns the vet id as int32
func getIdFromRequestContext(c *fiber.Ctx) (int32, string, error) {
	// Get the variable from the request context
	// Variable not found or not of type string
	vetIdStr, ok := c.Locals("userId").(string)
	if !ok {
		return 0, "Error", errors.New("error getting vet id")
	}

	// Convert the vetIdStr to int32
	vetId, err := strconv.Atoi(vetIdStr)
	if err != nil {
		return 0, "Invalid ID", err
	}

	return int32(vetId), "", nil
}
