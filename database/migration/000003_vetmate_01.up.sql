CREATE TABLE `vacunaciones` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `mascota_id` INT UNSIGNED,
  `tipo_vacuna_id` INT UNSIGNED,
  `vet_id` INT UNSIGNED,
  `direccion_sucursal_id` INT UNSIGNED,
  `fecha_aplicacion` date DEFAULT (now()),
  `laboratorio` varchar(150) NOT NULL,
  `lote_vacuna` varchar(255) NOT NULL,
  `peso` decimal(10,3) NOT NULL,
  `vacuna_fecha_caducidad` date NOT NULL,
  `prox_fecha_vacunacion` date NOT NULL
);

CREATE TABLE `vacunas` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `vet_id` INT UNSIGNED,
  `tipo_vacuna` varchar(200) NOT NULL
);

CREATE INDEX `vacunaciones_index_11` ON `vacunaciones` (`mascota_id`);

CREATE INDEX `vacunaciones_index_12` ON `vacunaciones` (`fecha_aplicacion`);

ALTER TABLE `vacunaciones` ADD FOREIGN KEY (`mascota_id`) REFERENCES `mascotas` (`id`);

ALTER TABLE `vacunaciones` ADD FOREIGN KEY (`tipo_vacuna_id`) REFERENCES `vacunas` (`id`);

ALTER TABLE `vacunaciones` ADD FOREIGN KEY (`direccion_sucursal_id`) REFERENCES `direccion_sucursales` (`id`);

ALTER TABLE `vacunas` ADD FOREIGN KEY (`vet_id`) REFERENCES `veterinarios` (`id`);
