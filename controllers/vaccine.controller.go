package controllers

import (
	"database/sql"

	db "github.com/HouseCham/VetMate/database/sql"
	"github.com/HouseCham/VetMate/validations"
	"github.com/gofiber/fiber/v2"
)

func InsertNewVaccineRecord(c *fiber.Ctx) error {
	var request db.Vacunacione

	// Parse request body from JSON to struct
	c.BodyParser(&request)
	// Trim input fields from request body
	purgeInputData(&request)

	// Validating user request parameters
	// if it is not valid, return 400 with error message
	if err := validations.ValidateRequest(&request, 1); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": responseMessages["invalidRequestBody"],
			"error":   err.Error(),
		})
	}

	insertNewVaccineChan := make(chan ErrorResponse)
	go func() {
		vetId, _, err := getIdFromRequestContext(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			insertNewVaccineChan <- ErrorResponse{
				Message: responseMessages["serverError"],
				Err:     err,
			}
		} else {
			params := db.InsertNewVaccineRecordParams{
				MascotaID:    request.MascotaID,
				TipoVacunaID: request.TipoVacunaID,
				VetID: sql.NullInt32{
					Int32: vetId,
					Valid: true,
				},
				DireccionSucursalID:  request.DireccionSucursalID,
				Laboratorio:          request.Laboratorio,
				LoteVacuna:           request.LoteVacuna,
				Peso:                 request.Peso,
				VacunaFechaCaducidad: request.VacunaFechaCaducidad,
				ProxFechaVacunacion:  request.ProxFechaVacunacion,
			}
			// Starting transaction
			tx, err := DB.Begin()
			// if error starting transaction
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				insertNewVaccineChan <- ErrorResponse{
					Message: responseMessages["serverError"],
					Err:     errorMessages["beginTX"],
				}
			} else {
				// implementing transaction in queries
				qtx := Queries.WithTx(tx)
				err = qtx.InsertNewVaccineRecord(c.Context(), params)
				// if error updating user
				if err != nil {
					c.Status(fiber.StatusInternalServerError)
					insertNewVaccineChan <- ErrorResponse{
						Message: responseMessages["serverError"],
						Err:     err,
					}
				} else {
					if err := tx.Commit(); err != nil {
						c.Status(fiber.StatusInternalServerError)
						insertNewVaccineChan <- ErrorResponse{
							Message: responseMessages["serverError"],
							Err:     errorMessages["commitTX"],
						}
					} else {
						insertNewVaccineChan <- ErrorResponse{
							Message: responseMessages["insertVaccineSuccess"],
							Err:     nil,
						}
					}
				}
			}
			defer tx.Rollback()
		}
		close(insertNewVaccineChan)
	}()

	// handling channel response
	if chanResponse := <-insertNewVaccineChan; chanResponse.Err != nil {
		return c.JSON(fiber.Map{
			"message": responseMessages["serverError"],
			"error":   chanResponse.Err.Error(),
		})
	} else {
		return c.JSON(fiber.Map{
			"message": responseMessages["insertVaccineSuccess"],
			"error":   nil,
		})
	}
}
