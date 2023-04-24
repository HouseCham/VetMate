-- name: InsertNewPet :exec
INSERT INTO mascotas (
    veterinaria_id,
    nombre,
    apellido_p,
    apellido_m,
    email,
    telefono,
    img_url,
    password_hash
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);