package controllers

import (
	"database/sql"
	"errors"

	"github.com/HouseCham/VetMate/config"
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/interfaces"
	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB
var Queries *db.Queries
var Config *config.Config

// ShareDbConnection is a function that shares the
// database connection to all controllers
// so that they can use the same connection
func ShareDbConnection(db *sql.DB) {
	DB = db
	Queries = createNewQuery()
}

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}

// createNewQuery is a function that creates a new
// query object for the database connection
func createNewQuery() *db.Queries {
	return db.New(DB)
}

// TrimInputFields is a function that trims all the
// input fields from the request body
// it is an interface function that is used by all
// the controllers
func trimInputFields(input interfaces.INewInsertParams) {
	input.Trim()
}

// CheckVetEmailAlreadyInUse is a function that checks
// if the email is already in use by checking the database
func checkVetEmailAlreadyInUse(email string, c *fiber.Ctx) (string, error) {
	emailExists, err := Queries.CheckVetEmailExists(c.Context(), email)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return "Error checking email", err
	} else if emailExists > 0 {
		return "Error", errors.New("email already in use")
	}
	return "", err
}