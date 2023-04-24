-- name: InsertNewUser :exec
INSERT INTO usuarios(
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?);