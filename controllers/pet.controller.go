package controllers

import (
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

// InsertNewPet inserts a new pet into the database
func InsertNewPetByUser(c *fiber.Ctx) error {
	var request db.Mascota
	var err error

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
		PropietarioID:   request.PropietarioID,
		RazaID:          request.RazaID,
		Descripcion:     request.Descripcion,
		Nombre:          request.Nombre,
		Sexo:            request.Sexo,
		Token:           util.RandomStringNum(10),
		ImgUrl:          request.ImgUrl,
		FechaNacimiento: request.FechaNacimiento,
	}
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
