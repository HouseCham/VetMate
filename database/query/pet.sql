-- name: InsertNewPet :exec
INSERT INTO mascotas (
    propietario_id,
    raza_id,
    token,
    descripcion,
    nombre,
    edad_aprox
) VALUES (?, ?, ?, ?, ?, ?);