package controllers

import (
	"database/sql"

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

func InsertNewVet(c *fiber.Ctx) error {
	var request db.InsertNewVetParams
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)
	// Hash password
	request.PasswordHash, err = util.HashPassword(request.PasswordHash)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	//! TODO: Check if email is already in use
	// Validate request body
	if isValid, err := validations.ValidateVet(request); !isValid {
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return Queries.InsertNewVet(c.Context(), request)
}