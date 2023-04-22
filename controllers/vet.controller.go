package controllers

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/interfaces"
	"github.com/HouseCham/VetMate/util"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

var DB *sql.DB
var Queries *db.Queries

//? ==================== COMMON ====================

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

// TrimInputFields is a function that trims all the
// input fields from the request body
// it is an interface function that is used by all
// the controllers
func trimInputFields(input interfaces.INewInsertParams) {
	input.Trim()
}

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

//? ==================== VET CONTROLLERS ====================

// InsertNewVet is a function that inserts a new vet
// to the database
func InsertNewVet(c *fiber.Ctx) error {
	var request db.InsertNewVetParams
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	trimInputFields(&request)

	// Check if email is already in use
	if message, err := checkVetEmailAlreadyInUse(request.Email, c); err != nil {
		c.Status(fiber.StatusConflict)
		return c.JSON(fiber.Map{
			"message": message,
			"error":   err.Error(),
		})
	}

	// Hash password
	// if error occurs, return 500
	request.PasswordHash, err = util.HashPassword(request.PasswordHash)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Error hashing password",
			"error":   err.Error(),
		})
	}

	// Validate request body
	// if not valid, return 400
	if isValid, err := validations.ValidateVet(request); !isValid {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
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
			"error":   err.Error(),
		})
	}
	// then, we need to convert the id to int32
	id32 := int32(id)
	// then, we need to get the vet info from the database
	// if error occurs, return 404
	mainInfo, err := Queries.GetVetMainInfoById(c.Context(), id32)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not get vet info",
			"error":   err.Error(),
		})
	}
	return c.JSON(mainInfo)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ? Implement Trim() function for LoginRequest struct
func (loginRequest *LoginRequest) Trim() {
	loginRequest.Email = strings.TrimSpace(loginRequest.Email)
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)
}

func LoginVet(c *fiber.Ctx) error {
	var request LoginRequest
	var err error

	// Parse request body from JSON to struct
	c.BodyParser(&request)

	// Trim input fields from request body
	trimInputFields(&request)

	// Validate request body
	// if not valid, return 400
	if isValid, err := validations.ValidateVetLogin(request.Email, request.Password); !isValid {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	// Get vet info from database
	vet, err := Queries.GetVetByEmail(c.Context(), request.Email)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Could not get vet info",
			"error":   err.Error(),
		})
	}
	// Compare password
	// if error occurs, return 500
	if err := util.CheckPassword(request.Password, vet.PasswordHash); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Wrong password",
			"error":   err.Error(),
		})
	}

	//! TODO: Generate JWT token

	return c.JSON(vet)
}
