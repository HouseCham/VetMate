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