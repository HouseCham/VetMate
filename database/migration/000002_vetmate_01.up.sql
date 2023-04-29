CREATE TABLE `usuarios` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `nombre` varchar(100) NOT NULL,
  `apellido_p` varchar(50) NOT NULL,
  `apellido_m` varchar(50) NOT NULL,
  `email` varchar(150) NOT NULL,
  `telefono` varchar(20),
  `password` varchar(255) NOT NULL,
  `email_validado` tinyint DEFAULT 0,
  `img_url` varchar(50),
  `calle` varchar(150) NOT NULL,
  `num_ext` varchar(10) NOT NULL,
  `num_int` varchar(10),
  `colonia` varchar(50) NOT NULL,
  `cp` varchar(10) NOT NULL,
  `ciudad` varchar(50) NOT NULL,
  `estado` varchar(50) NOT NULL,
  `pais` varchar(50) NOT NULL,
  `referencia` varchar(255),

  `fecha_registro` timestamp DEFAULT (NOW()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE TABLE `mascotas` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `propietario_id` INT UNSIGNED,
  `raza_id` INT UNSIGNED,
  `raza_comentario` varchar(100),
  `nombre` varchar(100),
  `edad` TINYINT,
  `proxima_fecha_vacunacion` timestamp,
  `proxima_fecha_desparasitacion` timestamp,

  `fecha_registro` timestamp DEFAULT (NOW()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE TABLE `familias` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `nombre` varchar(100)
);

CREATE TABLE `razas` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `familia_id` INT UNSIGNED,
  `nombre` varchar(100)
);

CREATE INDEX `nombre_completo_user` ON `usuarios` (`nombre`, `apellido_p`, `apellido_m`);

CREATE UNIQUE INDEX `usuarios_index_1` ON `usuarios` (`email`);

CREATE INDEX `nombre_mascota` ON `mascotas` (`nombre`);

CREATE INDEX `prox_fecha_vacuna` ON `mascotas` (`proxima_fecha_vacunacion`);

CREATE INDEX `prox_fecha_desparasito` ON `mascotas` (`proxima_fecha_desparasitacion`);

CREATE INDEX `familias_index_5` ON `familias` (`nombre`);

CREATE INDEX `razas_index_6` ON `razas` (`nombre`);

ALTER TABLE `mascotas` ADD FOREIGN KEY (`propietario_id`) REFERENCES `usuarios` (`id`);

ALTER TABLE `mascotas` ADD FOREIGN KEY (`raza_id`) REFERENCES `razas` (`id`);

ALTER TABLE `razas` ADD FOREIGN KEY (`familia_id`) REFERENCES `familias` (`id`);

-- Inserting default data during migration
INSERT INTO `familias` (`nombre`) VALUES('Cánido'), ('Félido');

-- Razas de perros
INSERT INTO `razas` (`familia_id`, `nombre`) VALUES(1, 'Otro'), (1, 'Shih tzu'),(1, 'Perdiguero de Burgos'),(1, 'Perdiguero Portugués'),(1, 'Cobrador de Pelo Rizado'),(1, 'Cobrador de Pelo Liso'),(1, 'Silky Terrier Australiano');

-- Razas de gatos