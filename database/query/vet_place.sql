-- name: CheckVetPlaceNameExists :one
SELECT COUNT(*)
FROM negocios
WHERE nombre_negocio = ?;

-- name: InsertNewVetPlace :exec
INSERT INTO negocios (nombre_negocio, token)
VALUES (?, ?);

-- name: UpdateVetPlace :exec
UPDATE negocios
SET nombre_negocio = ?
WHERE id = ?;

-- name: DeleteVetPlace :exec
UPDATE negocios
SET fecha_delete = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?;