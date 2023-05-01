-- name: InsertNewPet :exec
INSERT INTO mascotas (
    propietario_id,
    raza_id,
    descripcion,
    nombre,
    sexo,
    token,
    img_url,
    fecha_nacimiento
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: NewPetInValhalla :exec
UPDATE mascotas
SET fecha_muerte = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?;

-- name: DeletePet :exec
UPDATE mascotas
SET fecha_delete = DATE_SUB(NOW(), INTERVAL 6 HOUR)
WHERE id = ?;

-- name: GetPetMainInfo :one
SELECT mascotas.id, mascotas.nombre, razas.nombre as 'raza', CONCAT(usuarios.nombre, ' ', usuarios.apellido_p, ' ', usuarios.apellido_m) as 'propietario', descripcion, sexo, FLOOR(DATEDIFF(NOW(), fecha_nacimiento) / 365) as 'edad', mascotas.img_url, fecha_nacimiento, fecha_esterilizacion, ultima_fecha_desparasitacion, ultima_fecha_vacunacion 
FROM mascotas
INNER JOIN usuarios ON propietario_id = usuarios.id
INNER JOIN razas ON raza_id = razas.id
WHERE mascotas.fecha_delete IS NULL AND mascotas.id = ?;