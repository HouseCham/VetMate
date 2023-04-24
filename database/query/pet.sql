-- name: InsertNewPet :exec
INSERT INTO mascotas (
    propietario_id,
    raza_id,
    raza_comentario,
    nombre,
    edad
) VALUES (?, ?, ?, ?, ?);