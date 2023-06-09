// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: pet.sql

package db

import (
	"context"
	"database/sql"
)

const deletePet = `-- name: DeletePet :exec
UPDATE mascotas
SET fecha_delete = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?
`

func (q *Queries) DeletePet(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deletePet, id)
	return err
}

const getOwnerIdByPetId = `-- name: GetOwnerIdByPetId :one
SELECT propietario_id FROM mascotas WHERE id = ?
`

func (q *Queries) GetOwnerIdByPetId(ctx context.Context, id int32) (sql.NullInt32, error) {
	row := q.db.QueryRowContext(ctx, getOwnerIdByPetId, id)
	var propietario_id sql.NullInt32
	err := row.Scan(&propietario_id)
	return propietario_id, err
}

const getPetMainInfo = `-- name: GetPetMainInfo :one
SELECT mascotas.id, mascotas.nombre, razas.nombre as 'raza', CONCAT(usuarios.nombre, ' ', usuarios.apellido_p, ' ', usuarios.apellido_m) as 'propietario', descripcion, sexo, FLOOR(DATEDIFF(NOW(), fecha_nacimiento) / 365) as 'edad', mascotas.img_url, fecha_esterilizacion, ultima_fecha_desparasitacion, ultima_fecha_vacunacion 
FROM mascotas
INNER JOIN usuarios ON propietario_id = usuarios.id
INNER JOIN razas ON raza_id = razas.id
WHERE mascotas.fecha_delete IS NULL AND mascotas.id = ?
`

type GetPetMainInfoRow struct {
	ID                         int32          `json:"id"`
	Nombre                     sql.NullString `json:"nombre"`
	Raza                       sql.NullString `json:"raza"`
	Propietario                string         `json:"propietario"`
	Descripcion                sql.NullString `json:"descripcion"`
	Sexo                       string         `json:"sexo"`
	Edad                       int32          `json:"edad"`
	ImgUrl                     sql.NullString `json:"img_url"`
	FechaEsterilizacion        sql.NullTime   `json:"fecha_esterilizacion"`
	UltimaFechaDesparasitacion sql.NullTime   `json:"ultima_fecha_desparasitacion"`
	UltimaFechaVacunacion      sql.NullTime   `json:"ultima_fecha_vacunacion"`
}

func (q *Queries) GetPetMainInfo(ctx context.Context, id int32) (GetPetMainInfoRow, error) {
	row := q.db.QueryRowContext(ctx, getPetMainInfo, id)
	var i GetPetMainInfoRow
	err := row.Scan(
		&i.ID,
		&i.Nombre,
		&i.Raza,
		&i.Propietario,
		&i.Descripcion,
		&i.Sexo,
		&i.Edad,
		&i.ImgUrl,
		&i.FechaEsterilizacion,
		&i.UltimaFechaDesparasitacion,
		&i.UltimaFechaVacunacion,
	)
	return i, err
}

const insertNewPet = `-- name: InsertNewPet :exec
INSERT INTO mascotas (
    propietario_id,
    raza_id,
    descripcion,
    nombre,
    sexo,
    token,
    img_url,
    fecha_nacimiento
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertNewPetParams struct {
	PropietarioID   sql.NullInt32  `json:"propietario_id"`
	RazaID          sql.NullInt32  `json:"raza_id"`
	Descripcion     sql.NullString `json:"descripcion"`
	Nombre          sql.NullString `json:"nombre"`
	Sexo            string         `json:"sexo"`
	Token           string         `json:"token"`
	ImgUrl          sql.NullString `json:"img_url"`
	FechaNacimiento sql.NullTime   `json:"fecha_nacimiento"`
}

func (q *Queries) InsertNewPet(ctx context.Context, arg InsertNewPetParams) error {
	_, err := q.db.ExecContext(ctx, insertNewPet,
		arg.PropietarioID,
		arg.RazaID,
		arg.Descripcion,
		arg.Nombre,
		arg.Sexo,
		arg.Token,
		arg.ImgUrl,
		arg.FechaNacimiento,
	)
	return err
}

const newPetInValhalla = `-- name: NewPetInValhalla :exec
UPDATE mascotas
SET fecha_muerte = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?
`

func (q *Queries) NewPetInValhalla(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, newPetInValhalla, id)
	return err
}

const updatePet = `-- name: UpdatePet :exec
UPDATE mascotas
SET raza_id = ?, descripcion = ?, nombre = ?, sexo = ?, img_url = ?, fecha_nacimiento = ?, fecha_update = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?
`

type UpdatePetParams struct {
	RazaID          sql.NullInt32  `json:"raza_id"`
	Descripcion     sql.NullString `json:"descripcion"`
	Nombre          sql.NullString `json:"nombre"`
	Sexo            string         `json:"sexo"`
	ImgUrl          sql.NullString `json:"img_url"`
	FechaNacimiento sql.NullTime   `json:"fecha_nacimiento"`
	ID              int32          `json:"id"`
}

func (q *Queries) UpdatePet(ctx context.Context, arg UpdatePetParams) error {
	_, err := q.db.ExecContext(ctx, updatePet,
		arg.RazaID,
		arg.Descripcion,
		arg.Nombre,
		arg.Sexo,
		arg.ImgUrl,
		arg.FechaNacimiento,
		arg.ID,
	)
	return err
}
