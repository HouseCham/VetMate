package controllers

import (
	"database/sql"

	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

// InsertNewPet inserts a new pet into the database
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
			"message": "There is an error with the request",
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
			"message": "Error starting transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback()

	// implementing transaction in queries
	qtx := Queries.WithTx(tx)
	err = qtx.InsertNewPet(c.Context(), params)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Sorry, there was an error inserting pet info",
			"error":   err.Error(),
		})
	}

	return tx.Commit()
}

func InsertNewPetByVet(c *fiber.Ctx) error {
	panic("implement me") // TODO: Implement
}

func DeletePet(c *fiber.Ctx) error {
	panic("implement me") // TODO: Implement
}

func UpdatePet(c *fiber.Ctx) error {
	panic("implement me") // TODO: Implement
}
