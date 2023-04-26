-- name: GetVetMainInfoById :one
CALL getVetMainInfoById(?)

-- name: InsertNewVet :exec
INSERT INTO veterinarios (
    veterinaria_id,
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetVetByEmail :one
SELECT id, password_hash
FROM veterinarios
WHERE email = ?;

-- name: CheckVetEmailExists :one
SELECT COUNT(*)
FROM veterinarios
WHERE email = ?;

-- name: UpdateVet :exec
UPDATE veterinarios
SET nombre = ?, apellido_p = ?, apellido_m = ?, telefono = ?, img_url = ?, fecha_update = NOW()
WHERE id = ?;

-- name: DeleteVet :exec
UPDATE veterinarios
SET fecha_delete = NOW()
WHERE id = ?;