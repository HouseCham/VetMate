CREATE TABLE `veterinarios` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `sucursal_id` INT UNSIGNED,
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

CREATE TABLE `sucursales` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `nombre_sucursal` varchar(150) NOT NULL,
  `token` char(10) NOT NULL,

  `fecha_registro` timestamp DEFAULT (NOW()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE TABLE `direccion_sucursales` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `id_sucursal` INT UNSIGNED,
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

CREATE UNIQUE INDEX `sucursales_index_2` ON `sucursales` (`nombre_sucursal`);

CREATE INDEX `direccion_completa` ON `direccion_sucursales` (`calle`, `num_ext`, `colonia`);

ALTER TABLE `veterinarios` ADD FOREIGN KEY (`sucursal_id`) REFERENCES `sucursales` (`id`);

ALTER TABLE `direccion_sucursales` ADD FOREIGN KEY (`id_sucursal`) REFERENCES `sucursales` (`id`);