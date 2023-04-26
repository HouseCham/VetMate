CREATE PROCEDURE `getVetByEmail`(IN `vet_email` VARCHAR(150))
BEGIN
    SELECT `id`, `password_hash`
    FROM `veterinarios`
    WHERE `email` = `vet_email`;
END;

CREATE PROCEDURE `getVetMainInfoById`(IN `vet_id` INT UNSIGNED)
BEGIN
    SELECT `id`, `nombre`, `apellido_p`, `apellido_m`, `email`, `telefono`, `img_url`
    FROM `veterinarios`
    WHERE `id` = `vet_id`;
END;

CREATE PROCEDURE `insertNewVet`(IN `vet_name` VARCHAR(100), IN `vet_ap` VARCHAR(50), IN `vet_am` VARCHAR(50), IN `vet_email` VARCHAR(150), IN `vet_phone` VARCHAR(20), IN `vet_password_hash` VARCHAR(255))
BEGIN
    INSERT INTO `veterinarios` (
        `nombre`,
        `apellido_p`,
        `apellido_m`,
        `email`,
        `telefono`,
        `password_hash`
    ) VALUES (
        `vet_name`,
        `vet_ap`,
        `vet_am`,
        `vet_email`,
        `vet_phone`,
        `vet_password_hash`
    );
END;

CREATE PROCEDURE `checkVetEmailExists`(IN `vet_email` VARCHAR(150))
BEGIN
    SELECT COUNT(*)
    FROM `veterinarios`
    WHERE `email` = `vet_email`;
END;

CREATE PROCEDURE `updateVet`(IN `vet_id` INT UNSIGNED, IN `vet_name` VARCHAR(100), IN `vet_ap` VARCHAR(50), IN `vet_am` VARCHAR(50), IN `vet_phone` VARCHAR(20), IN `vet_img_url` VARCHAR(255))
BEGIN
    UPDATE `veterinarios`
    SET `nombre` = `vet_name`, `apellido_p` = `vet_ap`, `apellido_m` = `vet_am`, `telefono` = `vet_phone`, `img_url` = `vet_img_url`, `fecha_update` = NOW()
    WHERE `id` = `vet_id`;
END;

CREATE PROCEDURE `deleteVet`(IN `vet_id` INT UNSIGNED)
BEGIN
    UPDATE `veterinarios`
    SET `fecha_delete` = NOW()
    WHERE `id` = `vet_id`;
END;

CREATE PROCEDURE `getUserByEmail`(IN `user_email` VARCHAR(150))
BEGIN
    SELECT `id`, `password_hash`
    FROM `usuarios`
    WHERE `email` = `user_email`;
END;