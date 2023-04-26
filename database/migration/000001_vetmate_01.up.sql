CREATE TABLE `veterinarios` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `veterinaria_id` INT UNSIGNED,
  `nombre` varchar(100) NOT NULL,
  `apellido_p` varchar(50) NOT NULL,
  `apellido_m` varchar(50) NOT NULL,
  `email` varchar(150) NOT NULL UNIQUE,
  `telefono` varchar(20),
  `img_url` varchar(255) DEFAULT 'profile_404.png',
  `password_hash` varchar(255) NOT NULL,
  `email_validado` tinyint DEFAULT 0,

  `fecha_registro` timestamp DEFAULT (now()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE TABLE `negocios` (
  `id` INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  `nombre_negocio` varchar(150) NOT NULL,
  `guid` char(36) NOT NULL,

  `fecha_registro` timestamp DEFAULT (now()),
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
  `estado` varchar(50) NOT NULL,
  `pais` varchar(50) NOT NULL DEFAULT 'México',

  `fecha_registro` timestamp DEFAULT (now()),
  `fecha_update` timestamp,
  `fecha_delete` timestamp
);

CREATE INDEX `nombre_completo_vet` ON `veterinarios` (`nombre`, `apellido_p`, `apellido_m`);

CREATE UNIQUE INDEX `veterinarios_index_1` ON `veterinarios` (`email`);

CREATE UNIQUE INDEX `negocios_index_2` ON `negocios` (`nombre_negocio`);

CREATE INDEX `direccion_completa` ON `direccion_locales` (`calle`, `num_ext`, `colonia`);

ALTER TABLE `veterinarios` ADD FOREIGN KEY (`veterinaria_id`) REFERENCES `negocios` (`id`);

ALTER TABLE `direccion_locales` ADD FOREIGN KEY (`id_negocio`) REFERENCES `negocios` (`id`);


-- Inserting random data
INSERT INTO negocios(`nombre_negocio`, `guid`) VALUES('Marton Hospital Veterinario', 'd82ad76e-e456-11ed-b5ea-0242ac120002');
INSERT INTO direccion_locales(`id_negocio`,`calle`,`cp`,`num_ext`,`num_int`,`colonia`,`estado`) VALUES(1, 'Calle 1', '12345', '123', '123', 'Colonia 1', 'Estado 1');