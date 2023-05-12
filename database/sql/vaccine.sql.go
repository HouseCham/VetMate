// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: vaccine.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const getAllVaccinesFromPet = `-- name: GetAllVaccinesFromPet :many
SELECT vacunaciones.id, tipo_vacuna, lote_vacuna, fecha_aplicacion, nombre_sucursal, CONCAT(nombre + ' ' + apellido_p + ' ' + apellido_m) as 'nombre_veterinario'
FROM vacunaciones
INNER JOIN vacunas ON tipo_vacuna_id = vacunas.id
INNER JOIN sucursales ON direccion_sucursal_id = sucursales.id
INNER JOIN veterinarios ON sucursal_id = sucursales.id
WHERE mascota_id = ?
ORDER BY fecha_aplicacion DESC
`

type GetAllVaccinesFromPetRow struct {
	ID                int32        `json:"id"`
	TipoVacuna        string       `json:"tipo_vacuna"`
	LoteVacuna        string       `json:"lote_vacuna"`
	FechaAplicacion   sql.NullTime `json:"fecha_aplicacion"`
	NombreSucursal    string       `json:"nombre_sucursal"`
	NombreVeterinario string       `json:"nombre_veterinario"`
}

func (q *Queries) GetAllVaccinesFromPet(ctx context.Context, mascotaID sql.NullInt32) ([]GetAllVaccinesFromPetRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllVaccinesFromPet, mascotaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllVaccinesFromPetRow
	for rows.Next() {
		var i GetAllVaccinesFromPetRow
		if err := rows.Scan(
			&i.ID,
			&i.TipoVacuna,
			&i.LoteVacuna,
			&i.FechaAplicacion,
			&i.NombreSucursal,
			&i.NombreVeterinario,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertNewVaccineRecord = `-- name: InsertNewVaccineRecord :exec
INSERT INTO vacunaciones (
    mascota_id,
    tipo_vacuna_id,
    vet_id,
    direccion_sucursal_id,
    fecha_aplicacion,
    laboratorio,
    lote_vacuna,
    peso,
    vacuna_fecha_caducidad,
    prox_fecha_vacunacion
    ) VALUES (?, ?, ?, ?, DATE_SUB(NOW(), INTERVAL 6 HOUR), ?, ?, ?, ?, ?)
`

type InsertNewVaccineRecordParams struct {
	MascotaID            sql.NullInt32 `json:"mascota_id"`
	TipoVacunaID         sql.NullInt32 `json:"tipo_vacuna_id"`
	VetID                sql.NullInt32 `json:"vet_id"`
	DireccionSucursalID  sql.NullInt32 `json:"direccion_sucursal_id"`
	Laboratorio          string        `json:"laboratorio"`
	LoteVacuna           string        `json:"lote_vacuna"`
	Peso                 string        `json:"peso"`
	VacunaFechaCaducidad time.Time     `json:"vacuna_fecha_caducidad"`
	ProxFechaVacunacion  time.Time     `json:"prox_fecha_vacunacion"`
}

func (q *Queries) InsertNewVaccineRecord(ctx context.Context, arg InsertNewVaccineRecordParams) error {
	_, err := q.db.ExecContext(ctx, insertNewVaccineRecord,
		arg.MascotaID,
		arg.TipoVacunaID,
		arg.VetID,
		arg.DireccionSucursalID,
		arg.Laboratorio,
		arg.LoteVacuna,
		arg.Peso,
		arg.VacunaFechaCaducidad,
		arg.ProxFechaVacunacion,
	)
	return err
}
