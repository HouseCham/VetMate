-- name: CheckVetPlaceNameExists :one
SELECT COUNT(*)
FROM sucursales
WHERE nombre_sucursal = ?;

-- name: InsertNewVetPlace :exec
INSERT INTO sucursales (nombre_sucursal, token)
VALUES (?, ?);

-- name: UpdateVetPlace :exec
UPDATE sucursales
SET nombre_sucursal = ?
WHERE id = ?;

-- name: DeleteVetPlace :exec
UPDATE sucursales
SET fecha_delete = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?;