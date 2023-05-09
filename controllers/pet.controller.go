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
	// Getting the ownerId from the request context
	ownerId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

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
	}

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
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["beginTX"],
		})
	}
	defer tx.Rollback()

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.InsertNewPet(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["errorInsertingPet"],
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}

// InsertNewPet inserts a new pet into the database
// by a vet
func InsertNewPetByVet(c *fiber.Ctx) error {
	panic("implement me") // TODO: Implement
}

// DeletePet deletes a pet from the database
// updating fecha_delete to current date
func DeletePet(c *fiber.Ctx) error {
	// Getting the ownerId from the request context
	requestOwnerId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Getting the petId
	petId, err := getPetIdFromUri(c)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
		})
	}

	// Checking if the user is the owner of the pet
	if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   err.Error(),
		})
	} else if !isOwner {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": responseMessages["unauthorized"],
		})
	}

	// Starting transaction
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["beginTx"],
		})
	}
	defer tx.Rollback()

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.DeletePet(c.Context(), petId)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["deletePetError"],
			"error":   err.Error(),
		})
	}

	// Commiting transaction
	// if there is an error, return 500
	if err := tx.Commit(); err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["commitTx"],
		})
	}

	return c.JSON(fiber.Map{
		"message": responseMessages["deletePetSuccess"],
	})
}

// UpdatePet updates a pet in the database
// from a user
func UpdatePet(c *fiber.Ctx) error {
	// Getting the ownerId from the request context
	requestOwnerId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Getting the petId
	petId, err := getPetIdFromUri(c)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
		})
	}

	// Checking if the user is the owner of the pet
	if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   err.Error(),
		})
	} else if !isOwner {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": responseMessages["unauthorized"],
		})
	}

	// Parse request body from JSON to struct
	var request db.Mascota
	c.BodyParser(&request)

	purgeInputData(&request)

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 2); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}

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
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["beginTx"],
		})
	}
	defer tx.Rollback()

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.UpdatePet(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["updatePetError"],
			"error":   err.Error(),
		})
	}

	// Commiting transaction
	// if there is an error, return 500
	if err := tx.Commit(); err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["commitTx"],
		})
	}

	return c.JSON(fiber.Map{
		"message": responseMessages["updatePetSuccess"],
	})
}

// GetPet gets a pet from the database by a user
func GetPet(c *fiber.Ctx) error {
	// Getting the ownerId from the request context
	requestOwnerId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Getting the petId
	petId, err := getPetIdFromUri(c)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
		})
	}

	// Checking if the user is the owner of the pet
	if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   err.Error(),
		})
	} else if !isOwner {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": responseMessages["unauthorized"],
		})
	}

	// then, we need to get the vet info from the database
	// if error occurs, return 404
	mainInfo, err := Queries.GetPetMainInfo(c.Context(), petId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": responseMessages["petNotFound"],
			"error":   err.Error(),
		})
	}
	return c.JSON(mainInfo)
}

// NewWarriorInValhalla updates the 'fecha_muerte' column of a pet in the database to the current date. This is done when a pet dies and goes to Valhalla. This function can only be done by the owner of the pet
func NewWarriorInValhalla(c *fiber.Ctx) error {
	// Getting the ownerId from the request context
	requestOwnerId, message, err := getIdFromRequestContext(c)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Getting the petId
	petId, err := getPetIdFromUri(c)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
		})
	}

	// Checking if the user is the owner of the pet
	if isOwner, err := isUserOwner(requestOwnerId, petId, c); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   err.Error(),
		})
	} else if !isOwner {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": responseMessages["unauthorized"],
		})
	}

	// Starting transaction
	tx, err := DB.Begin()
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["beginTx"],
		})
	}
	defer tx.Rollback()

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.NewPetInValhalla(c.Context(), petId)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["petPassAwayError"],
			"error":   err.Error(),
		})
	}

	// Commiting transaction
	// if there is an error, return 500
	if err := tx.Commit(); err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   errorMessages["commitTx"],
		})
	}

	return c.JSON(fiber.Map{
		"message": responseMessages["petPassAwaySuccess"],
	})
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