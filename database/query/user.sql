-- name: InsertNewUser :exec
INSERT INTO usuarios(
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    password,
    calle,
    num_ext,
    num_int,
    colonia,
    cp,
    ciudad,
    estado,
    pais,
    referencia
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CheckUserEmailExists :one
SELECT COUNT(*)
FROM usuarios
WHERE email = ?;

-- name: GetUserByEmail :one
SELECT id, password
FROM usuarios
WHERE email = ?;

-- name: GetUserMainInfoById :one
SELECT nombre, apellido_p, apellido_m, email, telefono, img_url
FROM usuarios
WHERE id = ?;

-- name: UpdateUser :exec
UPDATE usuarios
SET nombre = ?, apellido_p = ?, apellido_m = ?,
telefono = ?, calle = ?, num_ext = ?, num_int = ?,
colonia = ?, cp = ?, ciudad = ?, estado = ?, pais = ?,
referencia = ?, fecha_update = NOW()
WHERE id = ?;

-- name: DeleteUser :exec
UPDATE usuarios
SET fecha_delete = NOW()
WHERE id = ?;