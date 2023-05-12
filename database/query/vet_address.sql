-- name: InsertBeginNewVetAddress :exec
INSERT INTO direccion_sucursales (id_negocio, calle, num_ext, num_int, colonia, cp, ciudad, estado, pais, referencia)
VALUES (LAST_INSERT_ID(), ?, ?, ?, ?, ?, ?, ?, ?, ?);