-- name: SelectRacesByFamily :many
SELECT id, nombre
FROM razas
WHERE familia_id = ?;