CREATE TABLE `veterinarios` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `veterinaria_id` INT UNSIGNED,
  `token` char(10) NOT NULL,
  `nombre` varchar(100) NOT NULL,
  `apellido_p` varchar(50) NOT NULL,
  `apellido_m` varchar(50) NOT NULL,
  `email` varchar(150) NOT NULL UNIQUE,
  `telefono` varchar(20),
  `img_url` varchar(255) DEFAULT 'profile_404.png',
  `password` varchar(255) NOT NULL,
  `email_validado` tinyint DEFAULT 0,

  `fecha_registro` timestamp DEFAULT (NOW()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
  -- NOW() -> DATE_SUB(NOW(), INTERVAL 6 HOUR)
);

CREATE TABLE `negocios` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `nombre_negocio` varchar(150) NOT NULL,
  `token` char(10) NOT NULL,

  `fecha_registro` timestamp DEFAULT (NOW()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE TABLE `direccion_locales` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `id_negocio` INT UNSIGNED,
  `calle` varchar(100) NOT NULL,
  `cp` varchar(10) NOT NULL,
  `num_ext` varchar(10) NOT NULL,
  `num_int` varchar(10),
  `colonia` varchar(50) NOT NULL,
  `ciudad` varchar(50) NOT NULL,
  `estado` varchar(50) NOT NULL,
  `pais` varchar(50) NOT NULL DEFAULT 'México',
  `referencia` varchar(255),

  `fecha_registro` timestamp DEFAULT (NOW()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE INDEX `nombre_completo_vet` ON `veterinarios` (`nombre`, `apellido_p`, `apellido_m`);

CREATE UNIQUE INDEX `veterinarios_index_1` ON `veterinarios` (`email`);

CREATE UNIQUE INDEX `negocios_index_2` ON `negocios` (`nombre_negocio`);

CREATE INDEX `direccion_completa` ON `direccion_locales` (`calle`, `num_ext`, `colonia`);

ALTER TABLE `veterinarios` ADD FOREIGN KEY (`veterinaria_id`) REFERENCES `negocios` (`id`);

ALTER TABLE `direccion_locales` ADD FOREIGN KEY (`id_negocio`) REFERENCES `negocios` (`id`);