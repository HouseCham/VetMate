package validations

import (
	"github.com/HouseCham/VetMate/config"
)

var Config *config.Config

// ErrorMessages is a map that contains the error messages
// that will be used in the validation process
var ErrorMessages = map[string]string{
	// English
	"required": "The %s is required",
	"email": "The %s must be a valid email address",
	"min": "The %s must be at least %s characters",
	"max": "The %s may not be greater than %s characters",
	"unique": "The %s has already been taken",
	"confirmed": "The %s confirmation does not match",
	"exists": "The %s does not exists",
	"date": "The %s must be a valid date",
	"date_format": "The %s does not match the format %s",
	"before": "The %s must be a date before %s",
	"after": "The %s must be a date after %s",
	"alpha": "The %s may only contain letters",
	"alpha_num": "The %s may only contain letters and numbers",
	"alpha_dash": "The %s may only contain letters, numbers, and dashes",
	"numeric": "The %s must be a number",

	// Spanish
	"requerido": "El campo %s es requerido",
	"correo": "El campo %s debe ser un correo válido",
	"minimo": "El campo %s debe tener al menos %d caracteres",
	"maximo": "El campo %s no debe tener mas de %d caracteres",
	"min_max": "El campo %s debe tener entre %d y %d caracteres",
	"unico": "El campo %s ya ha sido tomado",
	"confirmacion": "La confirmacion del campo %s no coincide",
	"existencia": "El campo %s no existe",
	"fecha": "El campo %s debe ser una fecha valida",
	"formato_fecha": "El campo %s no coincide con el formato %s",
	"antes": "El campo %s debe ser una fecha antes de %s",
	"despues": "El campo %s debe ser una fecha despues de %s",
	"alfabetico": "El campo %s solo puede contener letras",
	"alfanumerico": "El campo %s solo puede contener letras y numeros",
	"alfanumerico_guion": "El campo %s solo puede contener letras, numeros y guiones",
	"numerico": "El campo %s debe ser un número",
	"imagen": "El campo %s debe ser una imagen válida",
}

// ShareConfigFile is a function that shares the
// configuration setted up in main.go
func ShareConfigFile(config *config.Config) {
	Config = config
}