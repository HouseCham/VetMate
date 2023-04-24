-- name: InsertNewUser :exec
INSERT INTO usuarios(
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?);

-- name: CheckUserEmailExists :one
SELECT 1
FROM usuarios
WHERE email = ?;