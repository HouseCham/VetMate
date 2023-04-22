package controllers

import (
	"database/sql"
	"strconv"

	"github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB
var Queries *db.Queries

// ShareDbConnection is a function that shares the 
// database connection to all controllers
// so that they can use the same connection
func ShareDbConnection(db *sql.DB) {
	DB = db
	Queries = createNewQuery()
}

// createNewQuery is a function that creates a new
// query object for the database connection
func createNewQuery() *db.Queries {
	return db.New(DB)
}

// InsertNewVet is a function that inserts a new vet
// to the database
func InsertNewVet(c *fiber.Ctx) error {
	var request db.InsertNewVetParams
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)
	// Hash password
	// if error occurs, return 500
	request.PasswordHash, err = util.HashPassword(request.PasswordHash)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error hashing password",
			"error": err.Error(),
		})
	}

	//! TODO: Check if email is already in use
	// Validate request body
	// if not valid, return 400
	if isValid, err := validations.ValidateVet(request); !isValid {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
			"error": err.Error(),
		})
	}

	return Queries.InsertNewVet(c.Context(), request)
}

// GetVetById is a function that gets the vet info
// by the vet id from the url
func GetVetById(c *fiber.Ctx) error {
	// first, we need to get the id from the url
	// if error occurs, return 400
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid ID",
			"error": err.Error(),
		})
	}
	// then, we need to convert the id to int32
	id32 := int32(id)
	// then, we need to get the vet info from the database
	// if error occurs, return 404
	mainInfo, err :=  Queries.GetVetMainInfoById(c.Context(), id32)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not get vet info",
			"error": err.Error(),
		})
	}
	return c.JSON(mainInfo)
}