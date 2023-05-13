package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/HouseCham/VetMate/validations"
)

func (vaccine *Vacunacione) Trim() {
	vaccine.Laboratorio = strings.TrimSpace(vaccine.Laboratorio)
	vaccine.LoteVacuna = strings.TrimSpace(vaccine.LoteVacuna)
	vaccine.Peso = strings.TrimSpace(vaccine.Peso)
}

func (vaccine *Vacunacione) DeleteBlankFields() {
	vaccine.Laboratorio = strings.ReplaceAll(vaccine.Laboratorio, " ", "")
	vaccine.LoteVacuna = strings.ReplaceAll(vaccine.LoteVacuna, " ", "")
	vaccine.Peso = strings.ReplaceAll(vaccine.Peso, " ", "")
}

// ==================== VALIDATIONS ====================
func (vaccine *Vacunacione) ValidateNewRegister() error {
	// check if foreign keys are valid
	if !vaccine.MascotaID.Valid || vaccine.MascotaID.Int32 <= 0 {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "MascotaID")
	} else if !vaccine.TipoVacunaID.Valid || vaccine.TipoVacunaID.Int32 <= 0 {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "TipoVacunaID")
	} else if !vaccine.VetID.Valid || vaccine.VetID.Int32 <= 0 {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "VeterinarioID")
	} else if !vaccine.DireccionSucursalID.Valid || vaccine.DireccionSucursalID.Int32 <= 0 {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "VacunaID")
	}

	//check if any field is empty
	if vaccine.Laboratorio == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "Laboratorio")
	} else if vaccine.LoteVacuna == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "Lote de Vacuna")
	} else if vaccine.Peso == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "Peso")
	} else if vaccine.VacunaFechaCaducidad.String() == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "Fecha de Caducidad")
	} else if vaccine.ProxFechaVacunacion.String() == "" {
		return fmt.Errorf(validations.ErrorMessages["requerido"], "Proxima Fecha de Vacunacion")
	}

	// check fields length
	if len(vaccine.Laboratorio) > 150 {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "Laboratorio", 50)
	} else if len(vaccine.LoteVacuna) > 255 {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "Lote de Vacuna", 50)
	} else if len(vaccine.Peso) > 10 {
		return fmt.Errorf(validations.ErrorMessages["maximo"], "Peso", 10)
	}

	// check if date field format is valid
	if !isValidMySQLDate(vaccine.VacunaFechaCaducidad.String()) {
		return fmt.Errorf(validations.ErrorMessages["fecha"], "Fecha de Caducidad")
	} else if !isValidMySQLDate(vaccine.ProxFechaVacunacion.String()) {
		return fmt.Errorf(validations.ErrorMessages["fecha"], "Próxima Fecha de Vacunación")
	}

	if !isValidWeightValue(vaccine.Peso) {
		return errors.New(validations.ErrorMessages["peso"])
	} else if _, err := parseWeightValue(vaccine.Peso); err != nil {
		return errors.New(validations.ErrorMessages["peso"])
	}

	return nil
}

func (vaccine *Vacunacione) ValidateUpdate() error {
	panic("implement me")
}

func (vaccine *Vacunacione) ValidateLogin() error {
	panic("implement me")
}
