package controllers

import (
	"database/sql"
	"errors"
	"strconv"

	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

// InsertNewPet inserts a new pet into the database
// by a user
func InsertNewPetByUser(c *fiber.Ctx) error {
	var request db.Mascota
	// Parse request body from JSON to struct
	c.BodyParser(&request)
	purgeInputData(&request)
	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	} else {
		// Starting go routine
		newPetChan := make(chan error)
		go func() {
			// Getting the ownerId from the request context
			ownerId, message, err := getIdFromRequestContext(c)
			// if error getting the ownerId from request, return 500 with error message
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				newPetChan <- errors.New(message)
			} else {
				// mapping request into params
				params := db.InsertNewPetParams{
					PropietarioID:   sql.NullInt32{Int32: ownerId, Valid: true},
					RazaID:          request.RazaID,
					Descripcion:     request.Descripcion,
					Nombre:          request.Nombre,
					Sexo:            request.Sexo,
					Token:           util.RandomStringNum(10),
					ImgUrl:          request.ImgUrl,
					FechaNacimiento: request.FechaNacimiento,
				}
				// Starting transaction
				tx, err := DB.Begin()
				// if error starting transaction, return 500 with error message
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					newPetChan <- errorMessages["beginTX"]
				} else {
					// implementing transaction in queries
					qtx := Queries.WithTx(tx)
					err = qtx.InsertNewPet(c.Context(), params)
					// if error inserting new pet
					if err != nil {
						newPetChan <- err
					} else {
						newPetChan <- tx.Commit()
					}
				}
				defer tx.Rollback()
			}
			close(newPetChan)
		}()
		// Handling go routine response
		if err := <-newPetChan; err != nil {
			return c.JSON(fiber.Map{
				"message": responseMessages["insertPetError"],
				"error":   err.Error(),
			})
		} else {
			// Everything went well
			return c.JSON(fiber.Map{
				"message": responseMessages["insertPetSuccess"],
			})
		}
	}
}

// InsertNewPet inserts a new pet into the database
// by a vet
func InsertNewPetByVet(c *fiber.Ctx) error {
	panic("implement me") // TODO: Implement
}

// DeletePet deletes a pet from the database
// updating fecha_delete to current date
func DeletePet(c *fiber.Ctx) error {

	// Starting go routine
	deletePetChan := make(chan ErrorResponse)
	go func() {
		// Getting the ownerId from the request context
		requestOwnerId, message, err := getIdFromRequestContext(c)
		// if error getting the ownerId
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			deletePetChan <- ErrorResponse{
				Message: message,
				Err:     err,
			}
		} else {
			// Getting the petId
			petId, err := getPetIdFromUri(c)
			// if error getting the petId form uri
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				deletePetChan <- ErrorResponse{
					Message: responseMessages["invalidPetId"],
					Err:     err,
				}
			} else {
				// Checking if the user is the owner of the pet
				// return server error or not authorized if error
				if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
					c.Status(fiber.StatusInternalServerError)
					deletePetChan <- ErrorResponse{
						Message: responseMessages["serverError"],
						Err:     err,
					}
				} else if !isOwner {
					c.Status(fiber.StatusUnauthorized)
					deletePetChan <- ErrorResponse{
						Message: responseMessages["notAuthorized"],
						Err:     errorMessages["notOwner"],
					}
				} else {
					// Starting transaction
					tx, err := DB.Begin()
					// if error starting tx, return 500
					if err != nil {
						c.Status(fiber.StatusInternalServerError)
						deletePetChan <- ErrorResponse{
							Message: responseMessages["serverError"],
							Err:     errorMessages["beginTX"],
						}
					} else {
						// implementing transaction in queries
						qtx := Queries.WithTx(tx)
						err = qtx.DeletePet(c.Context(), petId)
						// if error deleting pet, return 500
						if err != nil {
							c.Status(fiber.StatusInternalServerError)
							deletePetChan <- ErrorResponse{
								Message: responseMessages["deletePetError"],
								Err:     err,
							}
						} else {
							// Commiting transaction
							// if there is an error, return 500
							if err := tx.Commit(); err != nil {
								c.Status(fiber.StatusInternalServerError)
								deletePetChan <- ErrorResponse{
									Message: responseMessages["serverError"],
									Err:     errorMessages["commitTX"],
								}
							} else {
								deletePetChan <- ErrorResponse{
									Message: "",
									Err:     nil,
								}
							}
						}
					}
					defer tx.Rollback()
				}
			}
		}
		close(deletePetChan)
	}()

	// Handling go routine response
	if serverResponse := <-deletePetChan; serverResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["deletePetError"],
			"error":   serverResponse.Err.Error(),
		})
	} else {
		// Everything went well
		return c.JSON(fiber.Map{
			"message": responseMessages["deletePetSuccess"],
		})
	}
}

// UpdatePet updates a pet in the database
// from a user
func UpdatePet(c *fiber.Ctx) error {
	// Parse request body from JSON to struct
	var request db.Mascota
	c.BodyParser(&request)

	purgeInputData(&request)

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	//* if bad request, return 400 with error message
	if err := validations.ValidateRequest(&request, 2); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	} else {
		// Starting goroutine
		updatePetChan := make(chan ErrorResponse)
		go func() {
			// Getting the ownerId from the request context
			requestOwnerId, message, err := getIdFromRequestContext(c)
			//* if error getting the ownerId
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				updatePetChan <- ErrorResponse{
					Message: message,
					Err:     err,
				}
			} else {
				// Getting the petId
				petId, err := getPetIdFromUri(c)
				//* if error getting the petId form uri
				if err != nil {
					c.Status(fiber.StatusBadRequest)
					updatePetChan <- ErrorResponse{
						Message: message,
						Err:     err,
					}
				} else {
					//* Checking if the user is the owner of the pet
					if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
						c.Status(fiber.StatusInternalServerError)
						updatePetChan <- ErrorResponse{
							Message: responseMessages["serverError"],
							Err:     err,
						}
					} else if !isOwner {
						c.Status(fiber.StatusUnauthorized)
						updatePetChan <- ErrorResponse{
							Message: responseMessages["unauthorized"],
							Err:     err,
						}
					} else {
						params := db.UpdatePetParams{
							RazaID:          request.RazaID,
							Descripcion:     request.Descripcion,
							Nombre:          request.Nombre,
							Sexo:            request.Sexo,
							ImgUrl:          request.ImgUrl,
							FechaNacimiento: request.FechaNacimiento,
							ID:              petId,
						}
						// Starting transaction
						tx, err := DB.Begin()
						//* if error starting tx, return 500
						if err != nil {
							updatePetChan <- ErrorResponse{
								Message: responseMessages["serverError"],
								Err:     errorMessages["beginTx"],
							}
						} else {
							// implementing transaction in queries
							qtx := Queries.WithTx(tx)
							err = qtx.UpdatePet(c.Context(), params)
							//* if error updating pet, return 500
							if err != nil {
								updatePetChan <- ErrorResponse{
									Message: responseMessages["updatePetError"],
									Err:     err,
								}
							} else {
								// Commiting transaction
								//* if there is an error, return 500
								if err := tx.Commit(); err != nil {
									updatePetChan <- ErrorResponse{
										Message: responseMessages["serverError"],
										Err:     errorMessages["commitTx"],
									}
								} else {
									updatePetChan <- ErrorResponse{
										Message: "",
										Err:     nil,
									}
								}
							}
						}
						defer tx.Rollback()
					}
				}
			}
			close(updatePetChan)
		}()

		// get the response from the channel
		if serverResponse := <-updatePetChan; serverResponse.Err != nil {
			return c.JSON(fiber.Map{
				"message": serverResponse.Message,
				"error":   serverResponse.Err.Error(),
			})
		} else {
			// Everything went well
			return c.JSON(fiber.Map{
				"message": responseMessages["updatePetSuccess"],
			})
		}
	}
}

// GetPet gets a pet from the database by a user
func GetPet(c *fiber.Ctx) error {
	// Starting goroutine
	getPetChan := make(chan HttpGetResponse)
	go func() {
		// Getting the ownerId from the request context
		requestOwnerId, message, err := getIdFromRequestContext(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			getPetChan <- HttpGetResponse{
				Message: message,
				Err:     err,
				Object:  db.GetPetMainInfoRow{},
			}
		} else {
			// Getting the petId
			petId, err := getPetIdFromUri(c)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				getPetChan <- HttpGetResponse{
					Message: responseMessages["invalidRequestBody"],
					Err:     nil,
					Object:  db.GetPetMainInfoRow{},
				}
			} else {
				// Checking if the user is the owner of the pet
				if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
					c.Status(fiber.StatusInternalServerError)
					getPetChan <- HttpGetResponse{
						Message: responseMessages["serverError"],
						Err:     err,
						Object:  db.GetPetMainInfoRow{},
					}
				} else if !isOwner {
					c.Status(fiber.StatusUnauthorized)
					getPetChan <- HttpGetResponse{
						Message: responseMessages["unauthorized"],
						Err:     nil,
						Object:  db.GetPetMainInfoRow{},
					}
				} else {
					// then, we need to get the vet info from the database
					// if error occurs, return 404
					mainInfo, err := Queries.GetPetMainInfo(c.Context(), petId)
					if err != nil {
						c.Status(fiber.StatusNotFound)
						getPetChan <- HttpGetResponse{
							Message: responseMessages["petNotFound"],
							Err:     err,
							Object:  db.GetPetMainInfoRow{},
						}
					} else {
						getPetChan <- HttpGetResponse{
							Message: "",
							Err:     nil,
							Object:  mainInfo,
						}
					}
				}
			}
		}
		close(getPetChan)
	}()

	// get the response from the channel
	serverResponse := <-getPetChan
	if serverResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": serverResponse.Message,
			"error":   serverResponse.Err.Error(),
		})
	} else {
		return c.JSON(serverResponse.Object)
	}
}

// NewWarriorInValhalla updates the 'fecha_muerte' column of a pet in the database to the current date. This is done when a pet dies and goes to Valhalla. This function can only be done by the owner of the pet
func NewWarriorInValhalla(c *fiber.Ctx) error {
	// Starting goroutine
	newWarriorChan := make(chan ErrorResponse)
	go func() {
		// Getting the ownerId from the request context
		requestOwnerId, message, err := getIdFromRequestContext(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			newWarriorChan <- ErrorResponse{
				Message: message,
				Err:     err,
			}
		} else {
			// Getting the petId
			petId, err := getPetIdFromUri(c)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				newWarriorChan <- ErrorResponse{
					Message: responseMessages["invalidRequestBody"],
					Err:     nil,
				}
			} else {
				// Checking if the user is the owner of the pet
				if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
					c.Status(fiber.StatusInternalServerError)
					newWarriorChan <- ErrorResponse{
						Message: responseMessages["serverError"],
						Err:     err,
					}
				} else if !isOwner {
					c.Status(fiber.StatusUnauthorized)
					newWarriorChan <- ErrorResponse{
						Message: responseMessages["unauthorized"],
						Err:     nil,
					}
				} else {
					// Starting transaction
					tx, err := DB.Begin()
					//? if error starting tx, return 500
					if err != nil {
						c.Status(fiber.StatusInternalServerError)
						newWarriorChan <- ErrorResponse{
							Message: responseMessages["serverError"],
							Err:     errorMessages["beginTx"],
						}
					} else {
						// implementing transaction in queries
						qtx := Queries.WithTx(tx)
						err = qtx.NewPetInValhalla(c.Context(), petId)
						//? if error updating pet, return 500
						if err != nil {
							c.Status(fiber.StatusInternalServerError)
							newWarriorChan <- ErrorResponse{
								Message: responseMessages["petPassAwayError"],
								Err:     err,
							}
						} else {
							// Commiting transaction
							//? if there is an error commiting tx, return 500
							if err := tx.Commit(); err != nil {
								c.Status(fiber.StatusInternalServerError)
								newWarriorChan <- ErrorResponse{
									Message: responseMessages["serverError"],
									Err:     errorMessages["commitTx"],
								}
							} else {
								newWarriorChan <- ErrorResponse{
									Message: "",
									Err:     nil,
								}
							}
						}
					}
					defer tx.Rollback()
				}
			}
		}
		close(newWarriorChan)
	}()

	// get the response from the channel
	if serverResponse := <-newWarriorChan; serverResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": serverResponse.Message,
			"error":   serverResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"message": responseMessages["petPassAwaySuccess"],
		})
	}
}

//? ================ functions used by controllers ================ ?//

// isUserOwner checks if the user that makes the request is the owner of the pet
func isUserOwner(requestOwnerId int32, petId int32, c *fiber.Ctx) (bool, error) {
	// Getting the ownerId from the request context
	ownerId, err := Queries.GetOwnerIdByPetId(c.Context(), petId)
	if err != nil {
		return false, err
	}
	// if the ownerId from the request context is null, equal to zero or is not the same as the ownerId from the database, return 401
	if !ownerId.Valid || ownerId.Int32 == 0 || ownerId.Int32 != requestOwnerId {
		return false, nil
	}
	return true, nil
}

// getPetIdFromUri gets the petId from the request uri
func getPetIdFromUri(c *fiber.Ctx) (int32, error) {
	// if id is empty, return 400
	id := c.Params("petId")
	if id == "" {
		return 0, errors.New("invalidRequestBody")
	}
	// converting petId to int32
	petId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(petId), nil
}
