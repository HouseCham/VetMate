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
	trimInputFields(&request)

	// Check if email is already in use
	// set false as second parameter in order to check for vet emails
	if message, err, status :=  checkEmailAlreadyInUse(request.Email, false, c); err != nil {
		c.Status(status)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Validate request body
	// if not valid, return 400
	if isValid, err := validations.ValidateVet(&request, 1); !isValid {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
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

	// Mapping request body to InsertNewVetParams struct
	params := db.InsertNewVetParams{
		Nombre:       request.Nombre,
		ApellidoP:    request.ApellidoP,
		ApellidoM:    request.ApellidoM,
		PasswordHash: request.PasswordHash,
		Email:        request.Email,
		Telefono:     request.Telefono,
		ImgUrl:       request.ImgUrl,
	}

	return Queries.InsertNewVet(c.Context(), params)
}
// GetVetById is a function that gets the vet info
// by the vet id from the url
func GetVetById(c *fiber.Ctx) error {
	// Get the variable from the request context
	// Variable not found or not of type string
	vetIdStr, ok := c.Locals("userId").(string)
	if !ok {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error getting vet id",
		})
	}

	// Convert the vetIdStr to int32
	vetId, err := strconv.Atoi(vetIdStr)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid ID",
			"error":   err.Error(),
		})
	}
	id32 := int32(vetId)

	// then, we need to get the vet info from the database
	// if error occurs, return 404
	mainInfo, err := Queries.GetVetMainInfoById(c.Context(), id32)
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
// by the vet id from the url
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
	trimInputFields(&request)

	// Validate request body
	// if not valid, return 400
	if isValid, err := validations.ValidateVet(&request, 2); !isValid {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Mapping request body to db.UpdateVetParams struct
	params := db.UpdateVetParams{
		ID:           request.ID,
		Nombre:       request.Nombre,
		ApellidoP:    request.ApellidoP,
		ApellidoM:    request.ApellidoM,
		Telefono:     request.Telefono,
		ImgUrl:       request.ImgUrl,
	}

	return Queries.UpdateVet(c.Context(), params)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ? Implement Trim() function for LoginRequest struct
func (loginRequest *LoginRequest) Trim() {
	loginRequest.Email = strings.TrimSpace(loginRequest.Email)
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)
}

// LoginVet is a function that logs in the vet
// by checking the email and password
func LoginVet(c *fiber.Ctx) error {
	var request LoginRequest
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	trimInputFields(&request)

	// Validate request body
	// if not valid, return 400
	if isValid, err := validations.ValidateVetLogin(request.Email, request.Password); !isValid {
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
	if err := util.CheckPassword(request.Password, vet.PasswordHash); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Wrong password",
			"error":   err.Error(),
		})
	}

	// Generating jwt
	// in case of error, returns 500 with error message
	tokenString, err := auth.GenerateJWT(vet.ID)
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