package controllers

import (
	"database/sql"
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
			"message": responseMessages[message],
			"error":   err.Error(),
		})
	}

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}

	// goroutine started to insert user into database
	insertChan := make(chan error)
	go func() {
		// Hash password
		// if error occurs, return 500
		request.Password, err = util.HashPassword(request.Password)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			insertChan <- errors.New("error encriptando contraseña")
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
			c.Status(fiber.StatusInternalServerError)
			insertChan <- errorMessages["beginTX"]
		}
		defer tx.Rollback()

		// implementing transaction in queries
		qtx := Queries.WithTx(tx)
		err = qtx.InsertNewUser(c.Context(), params)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			insertChan <- errorMessages["insertInfo"]
		}

		insertChan <- tx.Commit()
		close(insertChan)
	}()

	if err := <-insertChan; err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["errorInsertingUser"],
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": responseMessages["userInserted"],
	})
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

	// goroutine started to get user info
	getChan := make(chan db.GetUserMainInfoByIdRow)
	go func() {
		// then, we need to get the vet info from the database
		// if error occurs, return 404
		mainInfo, err := Queries.GetUserMainInfoById(c.Context(), userId)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			getChan <- db.GetUserMainInfoByIdRow{}
		}
		getChan <- mainInfo
		close(getChan)
	}()

	// getting user info from channel
	mainInfo := <-getChan
	if mainInfo == (db.GetUserMainInfoByIdRow{}) {
		return c.JSON(fiber.Map{
			"message": responseMessages["userNotFound"],
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
			"message": responseMessages["invalidRequestBody"],
		})
	}

	// goroutine started to login user
	loginChannel := make(chan LoginResponse)
	go func() {
		user, err := Queries.GetUserByEmail(c.Context(), request.Email)
		// DB error handling
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			// In case email not found
			if err == sql.ErrNoRows {
				loginChannel <- LoginResponse{
					Jwt: "",
					Err: errorMessages["wrongCredentials"],
				}
			} else {
				// In case of refused connection
				loginChannel <- LoginResponse{
					Jwt: "",
					Err: errorMessages["serverError"],
				}	
			}
		} else {
			// Compare dbPassword with request password
			// if error occurs, return 401
			if err := util.CheckPassword(request.Password, user.Password); err != nil {
				c.Status(fiber.StatusUnauthorized)
				loginChannel <- LoginResponse{
					Jwt: "",
					Err: errorMessages["wrongCredentials"],
				}
			} else {
				// Generating jwt
				tokenString, err := auth.GenerateJWT(user.ID, false)
				// in case of error, returns 500 with error message
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					loginChannel <- LoginResponse{
						Jwt: "",
						Err: errorMessages["generateJWT"],
					}
				} else {
					loginChannel <- LoginResponse{
						Jwt: tokenString,
						Err: nil,
					}
				}
			}
		}
		close(loginChannel)
	}()

	// handling channel response
	response := <-loginChannel
	if response.Err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["loginError"],
			"error":   response.Err.Error(),
		})
	}

	// if everything is ok, return jwt
	return c.JSON(fiber.Map{
		"message": responseMessages["loginSuccess"],
		"jwt": response.Jwt,
	})
}

// UpdateUser updates a user
func UpdateUser(c *fiber.Ctx) error {
	// Get the variable from the request context
	// Variable not found or not of type string
	userId, _, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": responseMessages["getIdError"],
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
			"message": responseMessages["invalidRequestBody"],
		})
	}

	// goroutine started to update user
	updateUserChan := make(chan error)
	go func() {
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
			updateUserChan <- errorMessages["transactionError"]
		}
		defer tx.Rollback()

		// implementing transaction in queries
		qtx := Queries.WithTx(tx)
		err = qtx.UpdateUser(c.Context(), params)
		if err != nil {
			updateUserChan <- errors.New("error al actualizar usuario")
		}

		updateUserChan <- tx.Commit()
		close(updateUserChan)
	}()

	if err := <-updateUserChan; err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["updateUserError"],
		})
	}

	return c.JSON(fiber.Map{
		"message": responseMessages["updateUserSuccess"],
	})
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

	// goroutine started to delete user
	deleteUserChan := make(chan error)
	go func() {
		// Starting transaction
		tx, err := DB.Begin()
		if err != nil {
			deleteUserChan <- errorMessages["beginTX"]
		}
		defer tx.Rollback()

		// implementing transaction in queries
		qtx := Queries.WithTx(tx)
		err = qtx.DeleteUser(c.Context(), userId)
		if err != nil {
			deleteUserChan <- errorMessages["deleteInfo"]
		}

		deleteUserChan <- tx.Commit()
		close(deleteUserChan)
	}()

	if err := <-deleteUserChan; err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["deleteUserError"],
		})
	}

	return c.JSON(fiber.Map{
		"message": responseMessages["deleteUserSuccess"],
	})
}
