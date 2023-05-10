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
	/* ========== Common ========== */
	"serverError":      errors.New("error en el servidor"),
	"getIdError":       errors.New("error al obtener id"),
	"beginTX":          errors.New("error al iniciar transacción"),
	"dbServerError":    errors.New("error con la base de datos"),
	
	/* ========== CRUD ========== */
	"updateInfo":       errors.New("error al actualizar información"),
	"insertInfo":       errors.New("error al insertar información"),
	"deleteInfo":       errors.New("error al eliminar información"),
	"getInfo":          errors.New("error al obtener información"),

	/* ========== Login-Authorization ========== */
	"generateJWT":      errors.New("error al generar token"),
	"hashPassword":     errors.New("error al encriptar contraseña"),
	"wrongPassword":    errors.New("contraseña incorrecta"),
	"wrongCredentials": errors.New("email o contraseña incorrectas"),
	"registerError":    errors.New("error al registrar usuario"),
	"notOwner":         errors.New("no eres el dueño de esta mascota"),
}

// responseMessages is a map that contains all the response messages
// that are going to be sent to the client
var responseMessages = map[string]string{
	/* ========== Common ========== */
	"serverError":        "Hubo un error en el servidor",
	"invalidRequestBody": "Cuerpo de la solicitud inválido",
	"notAuthorized":      "No autorizado",
	"commitTx":           "Error al confirmar transacción",

	/* ========== Login ========== */
	"emailInUse":   "El correo ya está en uso",
	"loginSuccess": "Sesión iniciada con éxito",
	"loginError":   "Hubo un error al iniciar sesión",

	/* ========== Register ========== */
	"registerSuccess": "Usuario registrado con éxito",
	"registerError":   "Hubo un error al registrar usuario",

	/* ========== VET ========== */
	// Insert
	"vetNotRegistered": "Hubo un error al registrar veterinario",
	"vetRegistered":    "Veterinario registrado con éxito",
	// Update
	"updateVetError":   "Hubo un error al actualizar veterinario",
	"updateVetSuccess": "Veterinario actualizado con éxito",
	// Delete
	"deleteVetError":   "Hubo un error al eliminar veterinario",
	"deleteVetSuccess": "Veterinario eliminado con éxito",
	// Get
	"vetNotFound": "Veterinario no encontrado",

	/* ========== USER ========== */
	// Insert
	"userInserted":       "Usuario registrado con éxito",
	"errorInsertingUser": "Hubo un error al registrar usuario",
	// Update
	"updateUserError":   "Hubo un error al actualizar usuario",
	"updateUserSuccess": "Usuario actualizado con éxito",
	// Delete
	"deleteUserError":   "Hubo un error al eliminar usuario",
	"deleteUserSuccess": "Usuario eliminado con éxito",
	// Get
	"userNotFound": "Usuario no encontrado",
	"userFound":    "Usuario encontrado",

	/* ========== PET ========== */
	// Insert
	"insertPetSuccess": "Mascota registrada con éxito",
	"insertPetError":   "Hubo un error al registrar mascota",
	// Update
	"updatePetError":   "Hubo un error al actualizar mascota",
	"updatePetSuccess": "Mascota actualizada con éxito",
	// Delete
	"deletePetError":   "Hubo un error al eliminar mascota",
	"deletePetSuccess": "Mascota eliminada con éxito",
	// Get
	"petNotFound":  "Mascota no encontrada",
	"petsNotFound": "Mascotas no encontradas",
	// Pass away
	"petPassAwayError":   "Hubo un error al registrar fallecimiento de mascota",
	"petPassAwaySuccess": "Lamentamos mucho tu pérdida, esperamos que estés bien",
}

// HttpGetResponse is the response
// for all the Get functions
// for GoRoutines
type HttpGetResponse struct {
	Message string      `json:"message"`
	Err     error       `json:"error"`
	Object  interface{} `json:"pet"`
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
	// Checking for email in use for vet or user table in database
	if isUserTable {
		emailExists, err = Queries.CheckUserEmailExists(c.Context(), email)
	} else {
		emailExists, err = Queries.CheckVetEmailExists(c.Context(), email)
	}
	// Handling server error or email already in use
	if err != nil {
		return responseMessages["serverError"], fiber.StatusInternalServerError, err
	} else if emailExists > 0 {
		return responseMessages["emailInUse"], fiber.StatusConflict, errorMessages["registerError"]
	}
	return "", 0, nil
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
