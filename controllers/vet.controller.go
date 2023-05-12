package controllers

import (
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
	// Parse request body from JSON to struct
	c.BodyParser(&request)
	// Trim input fields from request body
	purgeInputData(&request)
	// Validate request body
	// if not valid request, return 400
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": errorMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}
	// go routine for inserting vet start
	insertVetChan := make(chan ErrorResponse)
	go func() {
		// Check if email is already in use
		// set false as second parameter in order to check for vet emails
		//? if error occurs, return function state and error
		if message, status, err := checkEmailAlreadyInUse(request.Email, false, c); err != nil {
			c.Status(status)
			insertVetChan <- ErrorResponse{
				Message: message,
				Err:     err,
			}
		} else {
			// Hash password
			//? if error hashing password, return 500 and error
			request.Password, err = util.HashPassword(request.Password)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				insertVetChan <- ErrorResponse{
					Message: responseMessages["vetNotRegistered"],
					Err:     errorMessages["hashPassword"],
				}
			} else {
				// Mapping request body to InsertNewVetParams struct
				params := db.InsertNewVetParams{
					Nombre:    request.Nombre,
					ApellidoP: request.ApellidoP,
					ApellidoM: request.ApellidoM,
					Password:  request.Password,
					Email:     request.Email,
					Telefono:  request.Telefono,
					Token:     util.RandomStringNum(10),
				}
				// Starting transaction
				//? if error starting transaction, return 500 and error
				tx, err := DB.Begin()
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					insertVetChan <- ErrorResponse{
						Message: responseMessages["vetNotRegistered"],
						Err:     errorMessages["beginTX"],
					}
				} else {
					// getting queries with transaction
					qtx := Queries.WithTx(tx)
					err = qtx.InsertNewVet(c.Context(), params)
					//? if error inserting newVet
					if err != nil {
						c.Status(fiber.StatusInternalServerError)
						insertVetChan <- ErrorResponse{
							Message: responseMessages["vetNotRegistered"],
							Err:     errorMessages["insertInfo"],
						}
					} else {
						//? if error commiting tx
						if err := tx.Commit(); err != nil {
							insertVetChan <- ErrorResponse{
								Message: responseMessages["vetNotRegistered"],
								Err:     errorMessages["commitTX"],
							}
						} else {
							insertVetChan <- ErrorResponse{
								Message: responseMessages["vetRegistered"],
								Err:     nil,
							}
						}
					}
				}
				defer tx.Rollback()
			}
		}
		close(insertVetChan)
	}()

	//? handling channel response
	if chanResponse := <-insertVetChan; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"mensaje": chanResponse.Message,
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"mensaje": responseMessages["vetRegistered"],
		})
	}
}

// GetVetById is a function that gets the vet info
// by the vet id from the fiber context
func GetVetById(c *fiber.Ctx) error {
	getVetChan := make(chan HttpGetResponse)
	go func() {
		// Get the variable from the request context
		// Variable not found or not of type string
		vetId, message, err := getIdFromRequestContext(c)
		//? if error getting id from request context, return 500
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			getVetChan <- HttpGetResponse{
				Message: message,
				Err:     err,
				Object:  nil,
			}
		} else {
			// then, we need to get the vet info from the database
			mainInfo, err := Queries.GetVetMainInfoById(c.Context(), vetId)
			//? if error occurs, return 404
			if err != nil {
				c.Status(fiber.StatusNotFound)
				getVetChan <- HttpGetResponse{
					Message: responseMessages["vetNotFound"],
					Err:     err,
					Object:  nil,
				}
			} else {
				getVetChan <- HttpGetResponse{
					Message: responseMessages["vetFound"],
					Err:     nil,
					Object:  mainInfo,
				}
			}
		}
		close(getVetChan)
	}()

	//? handling channel response
	if chanResponse := <-getVetChan; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": chanResponse.Message,
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(chanResponse.Object)
	}
}

// UpdateVet is a function that updates the vet info
// by the vet id from the fiber context
func UpdateVet(c *fiber.Ctx) error {
	var request db.Veterinario

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	purgeInputData(&request)

	// Validate request body
	// if not valid, return 400
	if err := validations.ValidateRequest(&request, 2); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": errorMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}

	updateVetChan := make(chan ErrorResponse)
	go func() {
		// Get the vetId from the request context
		vetId, message, err := getIdFromRequestContext(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			updateVetChan <- ErrorResponse{
				Message: message,
				Err:     err,
			}
		} else {
			// Mapping request body to db.UpdateVetParams struct
			params := db.UpdateVetParams{
				ID:        vetId,
				Nombre:    request.Nombre,
				ApellidoP: request.ApellidoP,
				ApellidoM: request.ApellidoM,
				Telefono:  request.Telefono,
				ImgUrl:    request.ImgUrl,
			}
			// Starting transaction
			tx, err := DB.Begin()
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				updateVetChan <- ErrorResponse{
					Message: responseMessages["updateVetError"],
					Err:     errorMessages["beginTX"],
				}
			} else {
				// getting queries with transaction
				qtx := Queries.WithTx(tx)
				err = qtx.UpdateVet(c.Context(), params)
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					updateVetChan <- ErrorResponse{
						Message: responseMessages["updateVetError"],
						Err:     errorMessages["updateInfo"],
					}
				} else {
					if err := tx.Commit(); err != nil {
						c.Status(fiber.StatusInternalServerError)
						updateVetChan <- ErrorResponse{
							Message: responseMessages["updateVetError"],
							Err:     errorMessages["commitTX"],
						}
					} else {
						updateVetChan <- ErrorResponse{
							Message: responseMessages["updateVetSuccess"],
							Err:     nil,
						}
					}
				}
			}
			defer tx.Rollback()
		}
		close(updateVetChan)
	}()

	if chanResponse := <-updateVetChan; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"mensaje": chanResponse.Message,
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"mensaje": responseMessages["updateVetSuccess"],
		})
	}
}

// DeleteVet is a function that deletes the vet info
// by the vet id from fiber context
func DeleteVet(c *fiber.Ctx) error {
	// delete vet go routine start
	deleteVetChan := make(chan ErrorResponse)
	go func() {
		// Get the vetId from the request context
		vetId, message, err := getIdFromRequestContext(c)
		//? If variable not found or not of type string
		//? return 500 with error message
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			deleteVetChan <- ErrorResponse{
				Message: message,
				Err:     err,
			}
		} else {
			// Starting transaction
			tx, err := DB.Begin()
			//? if error starting transaction, return 500
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				deleteVetChan <- ErrorResponse{
					Message: responseMessages["deleteVetError"],
					Err:     errorMessages["beginTX"],
				}
			} else {
				// getting queries with transaction
				qtx := Queries.WithTx(tx)
				err = qtx.DeleteVet(c.Context(), vetId)
				//? if error deleting vet, return 500
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					deleteVetChan <- ErrorResponse{
						Message: responseMessages["deleteVetError"],
						Err:     errorMessages["deleteInfo"],
					}
				} else {
					//? if error on commit, return 500
					if err := tx.Commit(); err != nil {
						c.Status(fiber.StatusInternalServerError)
						deleteVetChan <- ErrorResponse{
							Message: responseMessages["deleteVetError"],
							Err:     errorMessages["commitTX"],
						}
					} else {
						deleteVetChan <- ErrorResponse{
							Message: responseMessages["deleteVetSuccess"],
							Err:     nil,
						}
					}
				}
			}
			defer tx.Rollback()
		}
		close(deleteVetChan)
	}()

	if chanResponse := <-deleteVetChan; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"mensaje": chanResponse.Message,
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"mensaje": responseMessages["deleteVetSuccess"],
		})
	}
}

// LoginVet is a function that logs in the vet
// by checking the email and password
func LoginVet(c *fiber.Ctx) error {
	var request db.Veterinario

	// Parse request body from JSON to struct
	c.BodyParser(&request)
	// Trim input fields from request body
	purgeInputData(&request)

	// Validate request body
	// if not valid, return 400
	if err := validations.ValidateRequest(&request, 3); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}

	// goroutine started to login vet
	loginVetChannel := make(chan LoginResponse)
	go func() {
		// Get vet info from database
		vet, err := Queries.GetVetByEmail(c.Context(), request.Email)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			loginVetChannel <- LoginResponse{
				Jwt: "",
				Err: errorMessages["getInfo"],
			}
		} else {
			// Compare password
			// if error occurs, return 500
			if err := util.CheckPassword(request.Password, vet.Password); err != nil {
				c.Status(fiber.StatusUnauthorized)
				loginVetChannel <- LoginResponse{
					Jwt: "",
					Err: errorMessages["wrongPassword"],
				}
			} else {
				// Generating jwt
				// in case of error, returns 500 with error message
				tokenString, err := auth.GenerateJWT(vet.ID, true)
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					loginVetChannel <- LoginResponse{
						Jwt: "",
						Err: errorMessages["generateJWT"],
					}
				} else {
					loginVetChannel <- LoginResponse{
						Jwt: tokenString,
						Err: nil,
					}
				}
			}
		}
		close(loginVetChannel)
	}()

	// Get response from goroutine
	if chanResponse := <-loginVetChannel; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["loginError"],
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"message": responseMessages["loginSuccess"],
			"jwt":     chanResponse.Jwt,
		})
	}
}
