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

// purgeInputData is a function that trims all the
// input fields and deletes blank spaces from the request body
// it is an interface function that is used by all
// the controllers
func purgeInputData(input interfaces.INewInsertParams) {
	input.Trim()
	input.DeleteBlankFields()
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
func checkEmailAlreadyInUse(email string, isUserTable bool, c *fiber.Ctx) (string, int, error) {
	var emailExists int64
	var err error

	if isUserTable {
		emailExists, err = Queries.CheckUserEmailExists(c.Context(), email)
	} else {
		emailExists, err = Queries.CheckVetEmailExists(c.Context(), email)
	}

	if err != nil {
		return "Hubo un error en el servidor", fiber.StatusInternalServerError, err
	} else if emailExists > 0 {
		return "conflicto", fiber.StatusConflict, errors.New("email ya usado por otro usuario")
	}
	return "Éxito", fiber.StatusOK, nil
}
