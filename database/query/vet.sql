-- name: GetVetMainInfoById :one
SELECT id, nombre, apellido_p, apellido_m, email, telefono, img_url
FROM veterinarios
WHERE id = ? AND fecha_delete IS NULL;

-- name: InsertNewVet :exec
INSERT INTO veterinarios (
    sucursal_id,
    token,
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetVetByEmail :one
SELECT id, password
FROM veterinarios
WHERE email = ? AND fecha_delete IS NULL;

-- name: CheckVetEmailExists :one
SELECT COUNT(*)
FROM veterinarios
WHERE email = ?;

-- name: UpdateVet :exec
UPDATE veterinarios
SET nombre = ?, apellido_p = ?, apellido_m = ?, telefono = ?, img_url = ?, fecha_update = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ? AND fecha_delete IS NULL;

-- name: DeleteVet :exec
UPDATE veterinarios
SET fecha_delete = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ? AND fecha_delete IS NULL;