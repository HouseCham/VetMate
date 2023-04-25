-- name: GetVetMainInfoById :one
SELECT id, nombre, apellido_p, apellido_m, email, telefono, img_url
FROM veterinarios
WHERE id = ?;

-- name: InsertNewVet :exec
INSERT INTO veterinarios (
    veterinaria_id,
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    img_url,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetVetByEmail :one
SELECT id, email, password_hash
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