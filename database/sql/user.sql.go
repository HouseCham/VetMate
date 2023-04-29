// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const checkUserEmailExists = `-- name: CheckUserEmailExists :one
SELECT COUNT(*)
FROM usuarios
WHERE email = ?
`

func (q *Queries) CheckUserEmailExists(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkUserEmailExists, email)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteUser = `-- name: DeleteUser :exec
UPDATE usuarios
SET fecha_delete = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ? AND fecha_delete IS NULL
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, password
FROM usuarios
WHERE email = ? AND fecha_delete IS NULL
`

type GetUserByEmailRow struct {
	ID       int32  `json:"id"`
	Password string `json:"password"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.ID, &i.Password)
	return i, err
}

const getUserMainInfoById = `-- name: GetUserMainInfoById :one
SELECT nombre, apellido_p, apellido_m, email, telefono, img_url
FROM usuarios
WHERE id = ? AND fecha_delete IS NULL
`

type GetUserMainInfoByIdRow struct {
	Nombre    string         `json:"nombre"`
	ApellidoP string         `json:"apellido_p"`
	ApellidoM string         `json:"apellido_m"`
	Email     string         `json:"email"`
	Telefono  sql.NullString `json:"telefono"`
	ImgUrl    sql.NullString `json:"img_url"`
}

func (q *Queries) GetUserMainInfoById(ctx context.Context, id int32) (GetUserMainInfoByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getUserMainInfoById, id)
	var i GetUserMainInfoByIdRow
	err := row.Scan(
		&i.Nombre,
		&i.ApellidoP,
		&i.ApellidoM,
		&i.Email,
		&i.Telefono,
		&i.ImgUrl,
	)
	return i, err
}

const insertNewUser = `-- name: InsertNewUser :exec
INSERT INTO usuarios(
    token,
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password,
    calle,
    num_ext,
    num_int,
    colonia,
    cp,
    ciudad,
    estado,
    pais,
    referencia
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertNewUserParams struct {
	Token      string         `json:"token"`
	Nombre     string         `json:"nombre"`
	ApellidoP  string         `json:"apellido_p"`
	ApellidoM  string         `json:"apellido_m"`
	Email      string         `json:"email"`
	Telefono   sql.NullString `json:"telefono"`
	Password   string         `json:"password"`
	Calle      string         `json:"calle"`
	NumExt     string         `json:"num_ext"`
	NumInt     sql.NullString `json:"num_int"`
	Colonia    string         `json:"colonia"`
	Cp         string         `json:"cp"`
	Ciudad     string         `json:"ciudad"`
	Estado     string         `json:"estado"`
	Pais       string         `json:"pais"`
	Referencia sql.NullString `json:"referencia"`
}

func (q *Queries) InsertNewUser(ctx context.Context, arg InsertNewUserParams) error {
	_, err := q.db.ExecContext(ctx, insertNewUser,
		arg.Token,
		arg.Nombre,
		arg.ApellidoP,
		arg.ApellidoM,
		arg.Email,
		arg.Telefono,
		arg.Password,
		arg.Calle,
		arg.NumExt,
		arg.NumInt,
		arg.Colonia,
		arg.Cp,
		arg.Ciudad,
		arg.Estado,
		arg.Pais,
		arg.Referencia,
	)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE usuarios
SET nombre = ?, apellido_p = ?, apellido_m = ?,
telefono = ?, calle = ?, num_ext = ?, num_int = ?,
colonia = ?, cp = ?, ciudad = ?, estado = ?, pais = ?,
referencia = ?, fecha_update = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ? AND fecha_delete IS NULL
`

type UpdateUserParams struct {
	Nombre     string         `json:"nombre"`
	ApellidoP  string         `json:"apellido_p"`
	ApellidoM  string         `json:"apellido_m"`
	Telefono   sql.NullString `json:"telefono"`
	Calle      string         `json:"calle"`
	NumExt     string         `json:"num_ext"`
	NumInt     sql.NullString `json:"num_int"`
	Colonia    string         `json:"colonia"`
	Cp         string         `json:"cp"`
	Ciudad     string         `json:"ciudad"`
	Estado     string         `json:"estado"`
	Pais       string         `json:"pais"`
	Referencia sql.NullString `json:"referencia"`
	ID         int32          `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Nombre,
		arg.ApellidoP,
		arg.ApellidoM,
		arg.Telefono,
		arg.Calle,
		arg.NumExt,
		arg.NumInt,
		arg.Colonia,
		arg.Cp,
		arg.Ciudad,
		arg.Estado,
		arg.Pais,
		arg.Referencia,
		arg.ID,
	)
	return err
}
