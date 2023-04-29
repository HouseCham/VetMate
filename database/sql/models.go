// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"database/sql"
)

type DireccionLocale struct {
	ID            int32          `json:"id"`
	IDNegocio     sql.NullInt32  `json:"id_negocio"`
	Calle         string         `json:"calle"`
	Cp            string         `json:"cp"`
	NumExt        string         `json:"num_ext"`
	NumInt        sql.NullString `json:"num_int"`
	Colonia       string         `json:"colonia"`
	Estado        string         `json:"estado"`
	Pais          string         `json:"pais"`
	FechaRegistro sql.NullTime   `json:"fecha_registro"`
	FechaUpdate   sql.NullTime   `json:"fecha_update"`
	FechaDelete   sql.NullTime   `json:"fecha_delete"`
}

type Familia struct {
	ID     int32          `json:"id"`
	Nombre sql.NullString `json:"nombre"`
}

type Mascota struct {
	ID                          int32          `json:"id"`
	PropietarioID               sql.NullInt32  `json:"propietario_id"`
	RazaID                      sql.NullInt32  `json:"raza_id"`
	Descripcion                 sql.NullString `json:"descripcion"`
	Nombre                      sql.NullString `json:"nombre"`
	EdadAprox                   sql.NullInt32  `json:"edad_aprox"`
	Sexo                        sql.NullString `json:"sexo"`
	Token                       string         `json:"token"`
	ImgUrl                      sql.NullString `json:"img_url"`
	FechaNacimiento             sql.NullTime   `json:"fecha_nacimiento"`
	FechaMuerte                 sql.NullTime   `json:"fecha_muerte"`
	FechaEsterilizacion         sql.NullTime   `json:"fecha_esterilizacion"`
	UltimaFechaDesparasitacion  sql.NullTime   `json:"ultima_fecha_desparasitacion"`
	UltimaFechaVacunacion       sql.NullTime   `json:"ultima_fecha_vacunacion"`
	ProximaFechaVacunacion      sql.NullTime   `json:"proxima_fecha_vacunacion"`
	ProximaFechaDesparasitacion sql.NullTime   `json:"proxima_fecha_desparasitacion"`
	FechaRegistro               sql.NullTime   `json:"fecha_registro"`
	FechaUpdate                 sql.NullTime   `json:"fecha_update"`
	FechaDelete                 sql.NullTime   `json:"fecha_delete"`
}

type Negocio struct {
	ID            int32        `json:"id"`
	NombreNegocio string       `json:"nombre_negocio"`
	Token         string       `json:"token"`
	FechaRegistro sql.NullTime `json:"fecha_registro"`
	FechaUpdate   sql.NullTime `json:"fecha_update"`
	FechaDelete   sql.NullTime `json:"fecha_delete"`
}

type Raza struct {
	ID        int32          `json:"id"`
	FamiliaID sql.NullInt32  `json:"familia_id"`
	Nombre    sql.NullString `json:"nombre"`
}

type Usuario struct {
	ID            int32          `json:"id"`
	Token         string         `json:"token"`
	Nombre        string         `json:"nombre"`
	ApellidoP     string         `json:"apellido_p"`
	ApellidoM     string         `json:"apellido_m"`
	Email         string         `json:"email"`
	Telefono      sql.NullString `json:"telefono"`
	Password      string         `json:"password"`
	EmailValidado sql.NullInt32  `json:"email_validado"`
	ImgUrl        sql.NullString `json:"img_url"`
	Calle         string         `json:"calle"`
	NumExt        string         `json:"num_ext"`
	NumInt        sql.NullString `json:"num_int"`
	Colonia       string         `json:"colonia"`
	Cp            string         `json:"cp"`
	Ciudad        string         `json:"ciudad"`
	Estado        string         `json:"estado"`
	Pais          string         `json:"pais"`
	Referencia    sql.NullString `json:"referencia"`
	FechaRegistro sql.NullTime   `json:"fecha_registro"`
	FechaUpdate   sql.NullTime   `json:"fecha_update"`
	FechaDelete   sql.NullTime   `json:"fecha_delete"`
}

type Veterinario struct {
	ID            int32          `json:"id"`
	VeterinariaID sql.NullInt32  `json:"veterinaria_id"`
	Token         string         `json:"token"`
	Nombre        string         `json:"nombre"`
	ApellidoP     string         `json:"apellido_p"`
	ApellidoM     string         `json:"apellido_m"`
	Email         string         `json:"email"`
	Telefono      sql.NullString `json:"telefono"`
	ImgUrl        sql.NullString `json:"img_url"`
	Password      string         `json:"password"`
	EmailValidado sql.NullInt32  `json:"email_validado"`
	FechaRegistro sql.NullTime   `json:"fecha_registro"`
	FechaUpdate   sql.NullTime   `json:"fecha_update"`
	FechaDelete   sql.NullTime   `json:"fecha_delete"`
}
