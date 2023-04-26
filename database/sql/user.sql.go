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

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, password_hash
FROM usuarios
WHERE email = ?
`

type GetUserByEmailRow struct {
	ID           int32  `json:"id"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(&i.ID, &i.PasswordHash)
	return i, err
}

const getUserMainInfoById = `-- name: GetUserMainInfoById :one
SELECT nombre, apellido_p, apellido_m, email, telefono, img_url
FROM usuarios
WHERE id = ?
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
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?)
`

type InsertNewUserParams struct {
	Nombre       string         `json:"nombre"`
	ApellidoP    string         `json:"apellido_p"`
	ApellidoM    string         `json:"apellido_m"`
	Email        string         `json:"email"`
	Telefono     sql.NullString `json:"telefono"`
	PasswordHash string         `json:"password_hash"`
}

func (q *Queries) InsertNewUser(ctx context.Context, arg InsertNewUserParams) error {
	_, err := q.db.ExecContext(ctx, insertNewUser,
		arg.Nombre,
		arg.ApellidoP,
		arg.ApellidoM,
		arg.Email,
		arg.Telefono,
		arg.PasswordHash,
	)
	return err
}
