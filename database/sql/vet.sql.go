// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: vet.sql

package db

import (
	"context"
	"database/sql"
)

const checkVetEmailExists = `-- name: CheckVetEmailExists :one
SELECT 1
FROM veterinarios
WHERE email = ?
`

func (q *Queries) CheckVetEmailExists(ctx context.Context, email string) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, checkVetEmailExists, email)
	var column_1 interface{}
	err := row.Scan(&column_1)
	return column_1, err
}

const getVetByEmail = `-- name: GetVetByEmail :one
SELECT id, email, password_hash
FROM veterinarios
WHERE email = ?
`

type GetVetByEmailRow struct {
	ID           int32  `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) GetVetByEmail(ctx context.Context, email string) (GetVetByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getVetByEmail, email)
	var i GetVetByEmailRow
	err := row.Scan(&i.ID, &i.Email, &i.PasswordHash)
	return i, err
}

const getVetMainInfoById = `-- name: GetVetMainInfoById :one
SELECT id, nombre, apellido_p, apellido_m, email, telefono, img_url
FROM veterinarios
WHERE id = ?
`

type GetVetMainInfoByIdRow struct {
	ID        int32          `json:"id"`
	Nombre    string         `json:"nombre"`
	ApellidoP string         `json:"apellido_p"`
	ApellidoM string         `json:"apellido_m"`
	Email     string         `json:"email"`
	Telefono  sql.NullString `json:"telefono"`
	ImgUrl    sql.NullString `json:"img_url"`
}

func (q *Queries) GetVetMainInfoById(ctx context.Context, id int32) (GetVetMainInfoByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getVetMainInfoById, id)
	var i GetVetMainInfoByIdRow
	err := row.Scan(
		&i.ID,
		&i.Nombre,
		&i.ApellidoP,
		&i.ApellidoM,
		&i.Email,
		&i.Telefono,
		&i.ImgUrl,
	)
	return i, err
}

const insertNewVet = `-- name: InsertNewVet :exec
INSERT INTO veterinarios (
    veterinaria_id,
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    img_url,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertNewVetParams struct {
	VeterinariaID sql.NullInt32  `json:"veterinaria_id"`
	Nombre        string         `json:"nombre"`
	ApellidoP     string         `json:"apellido_p"`
	ApellidoM     string         `json:"apellido_m"`
	Email         string         `json:"email"`
	Telefono      sql.NullString `json:"telefono"`
	ImgUrl        sql.NullString `json:"img_url"`
	PasswordHash  string         `json:"password_hash"`
}

func (q *Queries) InsertNewVet(ctx context.Context, arg InsertNewVetParams) error {
	_, err := q.db.ExecContext(ctx, insertNewVet,
		arg.VeterinariaID,
		arg.Nombre,
		arg.ApellidoP,
		arg.ApellidoM,
		arg.Email,
		arg.Telefono,
		arg.ImgUrl,
		arg.PasswordHash,
	)
	return err
}
