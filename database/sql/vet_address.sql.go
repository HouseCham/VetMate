// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: vet_address.sql

package db

import (
	"context"
	"database/sql"
)

const insertBeginNewVetAddress = `-- name: InsertBeginNewVetAddress :exec
INSERT INTO direccion_locales (id_negocio, calle, num_ext, num_int, colonia, cp, ciudad, estado, pais, referencia)
VALUES (LAST_INSERT_ID(), ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type InsertBeginNewVetAddressParams struct {
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

func (q *Queries) InsertBeginNewVetAddress(ctx context.Context, arg InsertBeginNewVetAddressParams) error {
	_, err := q.db.ExecContext(ctx, insertBeginNewVetAddress,
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
