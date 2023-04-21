package controllers

import (
	"database/sql"

	"github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/util"
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

	c.BodyParser(&request)
	request.PasswordHash, err = util.HashPassword(request.PasswordHash)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	return Queries.InsertNewVet(c.Context(), request)
}
