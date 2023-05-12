-- name: InsertNewVaccineRecord :exec
INSERT INTO vacunaciones (
    mascota_id,
    tipo_vacuna_id,
    vet_id,
    direccion_sucursal_id,
    fecha_aplicacion,
    laboratorio,
    lote_vacuna,
    peso,
    vacuna_fecha_caducidad,
    prox_fecha_vacunacion
    ) VALUES (?, ?, ?, ?, DATE_SUB(NOW(), INTERVAL 6 HOUR), ?, ?, ?, ?, ?);

-- name: GetAllVaccinesFromPet :many
SELECT vacunaciones.id, tipo_vacuna, lote_vacuna, fecha_aplicacion, nombre_sucursal, CONCAT(nombre + ' ' + apellido_p + ' ' + apellido_m) as 'nombre_veterinario'
FROM vacunaciones
INNER JOIN vacunas ON tipo_vacuna_id = vacunas.id
INNER JOIN sucursales ON direccion_sucursal_id = sucursales.id
INNER JOIN veterinarios ON sucursal_id = sucursales.id
WHERE mascota_id = ?
ORDER BY fecha_aplicacion DESC;