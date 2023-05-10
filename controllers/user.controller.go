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

	// goroutine started to insert user into database
	insertChan := make(chan ErrorResponse)
	go func() {
		// Check if email is already in use
		if message, status, err := checkEmailAlreadyInUse(request.Email, true, c); err != nil {
			c.Status(status)
			insertChan <- ErrorResponse{
				Message: message,
				Err:     err,
			}
		} else {
			// Hash password
			request.Password, err = util.HashPassword(request.Password)
			// if error occurs, return 500
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				insertChan <- ErrorResponse{
					Message: responseMessages["registerError"],
					Err:     err,
				}
			} else {
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
				// if error occurs starting tx, return 500
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					insertChan <- ErrorResponse{
						Message: responseMessages["registerError"],
						Err:     errorMessages["beginTX"],
					}
				} else {
					// implementing transaction in queries
					qtx := Queries.WithTx(tx)
					err = qtx.InsertNewUser(c.Context(), params)
					// if error occurs commiting tx, return 500
					if err != nil {
						c.Status(fiber.StatusInternalServerError)
						insertChan <- ErrorResponse{
							Message: responseMessages["registerError"],
							Err:     errorMessages["insertInfo"],
						}
					} else {
						insertChan <- ErrorResponse{
							Message: responseMessages["commitTx"],
							Err:     tx.Commit(),
						}
					}

				}
				defer tx.Rollback()
			}
		}
		close(insertChan)
	}()

	// handling response from goroutine
	if chanResponse := <-insertChan; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": chanResponse.Message,
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"message": responseMessages["userInserted"],
		})
	}
}

// GetUserByEmail gets a user by email
// if user does not exist, return 404
func GetUserById(c *fiber.Ctx) error {
	// goroutine started to get user info
	getChan := make(chan HttpGetResponse)
	go func() {
		// Get the variable from the request context
		// Variable not found or not of type string
		userId, message, err := getIdFromRequestContext(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			getChan <- HttpGetResponse{
				Message: message,
				Err:     err,
				Object:  db.GetUserMainInfoByIdRow{},
			}
		} else {
			// then, we need to get the vet info from the database
			mainInfo, err := Queries.GetUserMainInfoById(c.Context(), userId)
			// handling db select error
			if err != nil {
				// if no user is returned
				if errors.Is(err, sql.ErrNoRows) {
					c.Status(fiber.StatusNotFound)
					getChan <- HttpGetResponse{
						Message: responseMessages["userNotFound"],
						Err:     err,
						Object:  db.GetUserMainInfoByIdRow{},
					}
				} else {
					// in case there is a server error
					c.Status(fiber.StatusInternalServerError)
					getChan <- HttpGetResponse{
						Message: responseMessages["serverError"],
						Err:     err,
						Object:  db.GetUserMainInfoByIdRow{},
					}
				}
			} else {
				//returning user found
				getChan <- HttpGetResponse{
					Message: responseMessages["userFound"],
					Err:     nil,
					Object:  mainInfo,
				}
			}
		}
		close(getChan)
	}()

	// getting user info from channel
	chanResponse := <-getChan
	if chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": chanResponse.Message,
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"message": chanResponse.Message,
			"user":    chanResponse.Object,
		})
	}
}

// LoginResponse is a struct that contains
// the jwt and error returned from login functions
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
	chanResponse := <-loginChannel
	if chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["loginError"],
			"error":   chanResponse.Err.Error(),
		})
	} else {
		// if everything is ok, return jwt
		return c.JSON(fiber.Map{
			"message": responseMessages["loginSuccess"],
			"jwt":     chanResponse.Jwt,
		})
	}
}

// UpdateUser updates a user
func UpdateUser(c *fiber.Ctx) error {
	var request db.Usuario
	// Parse request body from JSON to struct
	c.BodyParser(&request)
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
		// Get the variable from the request context
		// Variable not found or not of type string
		userId, _, err := getIdFromRequestContext(c)
		// if error getting id from request
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			updateUserChan <- errorMessages["getIdError"]
		} else {
			// mapping request body to update user database params
			params := db.UpdateUserParams{
				ID:         userId,
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
			// if error starting transaction
			if err != nil {
				updateUserChan <- errorMessages["beginTX"]
			} else {
				// implementing transaction in queries
				qtx := Queries.WithTx(tx)
				err = qtx.UpdateUser(c.Context(), params)
				// if error updating user
				if err != nil {
					updateUserChan <- errorMessages["dbServerError"]
				} else {
					updateUserChan <- tx.Commit()
				}
			}
			defer tx.Rollback()
		}
		close(updateUserChan)
	}()

	// handling channel response in case of error
	if chanResponseErr := <-updateUserChan; chanResponseErr != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["updateUserError"],
			"error": chanResponseErr.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"message": responseMessages["updateUserSuccess"],
		})
	}
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
