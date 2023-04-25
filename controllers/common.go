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

type IsEmailUsedChan struct {
	Status  int
	Message string
	Err     error
}

// CheckVetEmailAlreadyInUse is a function that checks
// if the email is already in use by checking the database.
// if isUserTable is true, then it means we are trying to check if an user email already exists in database
// otherwise, if it is false, we are trying to check for a vet's email.
func checkEmailAlreadyInUse(email string, isUserTable bool, c *fiber.Ctx) (string, error, int) {
	var emailExists int64
	var err error

	if isUserTable {
		emailExists, err = Queries.CheckUserEmailExists(c.Context(), email)
	} else {
		emailExists, err = Queries.CheckVetEmailExists(c.Context(), email)
	}

	if err != nil {
		return "there was an error on the server, try later", err, fiber.StatusInternalServerError
	} else if emailExists > 0 {
		return "conflict", errors.New("email already in use"), fiber.StatusConflict
	}
	return "success", nil, fiber.StatusOK
}
