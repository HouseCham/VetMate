package controllers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/HouseCham/VetMate/config"
	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/interfaces"
	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB
var Queries *db.Queries
var Config *config.Config

// errorMessages is a map that contains all the error messages
// that are going to be sent to the client
var errorMessages = map[string]error{
	"beginTX": errors.New("error al iniciar transacción"),
	"updateInfo": errors.New("error al actualizar información"),
	"insertInfo": errors.New("error al insertar información"),
	"hashPassword": errors.New("error al encriptar contraseña"),
	"deleteInfo": errors.New("error al eliminar información"),
	"getInfo": errors.New("error al obtener información"),
	"wrongPassword": errors.New("contraseña incorrecta"),
	"generateJWT": errors.New("error al generar token"),
}

// responseMessages is a map that contains all the response messages
// that are going to be sent to the client
var responseMessages = map[string]string{
	"vetNotRegistered": "Hubo un error al registrar veterinario",
	"vetRegistered": "Veterinario registrado con éxito",
	"invalidRequestBody": "Cuerpo de la solicitud inválido",
	"noVetFound": "No se encontró veterinario",
	"updateVetError": "Hubo un error al actualizar veterinario",
	"updateVetSuccess": "Veterinario actualizado con éxito",
	"deleteVetError": "Hubo un error al eliminar veterinario",
	"deleteVetSuccess": "Veterinario eliminado con éxito",
	"loginError": "Hubo un error al iniciar sesión",
	"loginSuccess": "Sesión iniciada con éxito",
}

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

// getIdFromRequestContext is a function that gets the vet id
// from the request context -> c.Locals("userId")
// returns the vet id as int32
func getIdFromRequestContext(c *fiber.Ctx) (int32, string, error) {
	// Get the variable from the request context
	// Variable not found or not of type string
	vetIdStr, ok := c.Locals("userId").(string)
	if !ok {
		return 0, "Error", errors.New("error getting vet id")
	}

	// Convert the vetIdStr to int32
	vetId, err := strconv.Atoi(vetIdStr)
	if err != nil {
		return 0, "Invalid ID", err
	}

	return int32(vetId), "", nil
}