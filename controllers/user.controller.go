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
	if message, status, err := checkEmailAlreadyInUse(request.Email, true, c); err != nil {
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
		Nombre:     request.Nombre,
		ApellidoP:  request.ApellidoP,
		ApellidoM:  request.ApellidoM,
		Email:      request.Email,
		Telefono:   request.Telefono,
		Password:   request.Password,
		Calle:      request.Calle,
		Colonia:    request.Colonia,
		Ciudad:     request.Ciudad,
		Estado:     request.Estado,
		Cp:         request.Cp,
		Pais:       request.Pais,
		NumExt:     request.NumExt,
		NumInt:     request.NumInt,
		Referencia: request.Referencia,
		Token:      util.RandomStringNum(10),
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

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.InsertNewUser(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Sorry, there was an error inserting user info",
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}

// GetUserByEmail gets a user by email
// if user does not exist, return 404
func GetUserById(c *fiber.Ctx) error {
	// Get the variable from the request context
	// Variable not found or not of type string
	userId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// then, we need to get the vet info from the database
	// if error occurs, return 404
	mainInfo, err := Queries.GetUserMainInfoById(c.Context(), userId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not get user info",
			"error":   err.Error(),
		})
	}
	return c.JSON(mainInfo)
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
	Err error  `json:"error"`
}

// LoginUser is a function that logs in a user
// matching email and password
func LoginUser(c *fiber.Ctx) error {
	var request db.Usuario

	c.BodyParser(&request)
	purgeInputData(&request)

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 3); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"mensaje": "Hubo un error al iniciar sesión",
			"error":   "solicitud incorrecta",
		})
	}

	loginChannel := make(chan LoginResponse)
	go func() {

		user, err := Queries.GetUserByEmail(c.Context(), request.Email)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			loginChannel <- LoginResponse{
				Jwt: "",
				Err: errors.New("usuario o contraseña incorrecta"),
			}
		}

		// Compare password
		// if error occurs, return 401
		if err := util.CheckPassword(request.Password, user.Password); err != nil {
			c.Status(fiber.StatusUnauthorized)
			loginChannel <- LoginResponse{
				Jwt: "",
				Err: errors.New("usuario o contraseña incorrecta"),
			}
		}

		// Generating jwt
		// in case of error, returns 500 with error message
		tokenString, err := auth.GenerateJWT(user.ID, false)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			loginChannel <- LoginResponse{
				Jwt: "",
				Err: err,
			}
		}

		loginChannel <- LoginResponse{
			Jwt: tokenString,
			Err: nil,
		}
	}()

	response := <-loginChannel
	if response.Err != nil {
		return c.JSON(fiber.Map{
			"mensaje": "Hubo un error al iniciar sesión",
			"error":   response.Err.Error(),
		})
	}
	return c.JSON(response)
}

// UpdateUser updates a user
func UpdateUser(c *fiber.Ctx) error {
	// Get the variable from the request context
	// Variable not found or not of type string
	userId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	var request db.Usuario

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	request.ID = userId

	// Trim() and deleting blank spaces from request body
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

	// mapping request body to update user database params
	params := db.UpdateUserParams{
		ID:         request.ID,
		Nombre:     request.Nombre,
		ApellidoP:  request.ApellidoP,
		ApellidoM:  request.ApellidoM,
		Telefono:   request.Telefono,
		Calle:      request.Calle,
		Colonia:    request.Colonia,
		Ciudad:     request.Ciudad,
		Estado:     request.Estado,
		Cp:         request.Cp,
		Pais:       request.Pais,
		NumExt:     request.NumExt,
		NumInt:     request.NumInt,
		Referencia: request.Referencia,
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

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.UpdateUser(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error updating vet",
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}

// DeleteUser deletes a user
// by updating the fecha_delete field to current date
func DeleteUser(c *fiber.Ctx) error {
	// Get the variable from the request context
	// Variable not found or not of type string
	userId, message, err := getIdFromRequestContext(c)
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

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.DeleteUser(c.Context(), userId)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error deleting user",
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}
