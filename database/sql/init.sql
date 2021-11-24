CREATE DATABASE IF NOT EXISTS golang;


CREATE USER 'gowebserver'@'%' IDENTIFIED BY 'gopassword';
GRANT all ON golang.* TO 'gowebserver'@'%';

use golang;

grant all on *.* to gowebserver@'%' identified by 'gopassword' with grant option;

flush privileges;


-- CREATE TABLE IF NOT EXISTS `golang`.`golang_user` (
--   `user_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
--   `user_mail` VARCHAR(255) NULL,
--   `user_name` VARCHAR(255) NULL,
--   `user_password` VARCHAR(255) NULL,
--   `user_nickname` VARCHAR(255) NULL,
--   `user_registerDate` DATETIME NULL,
--   `user_lastLoginTime` DATETIME NULL,
--   `user_salt` VARCHAR(255) NULL,
--   PRIMARY KEY (`user_id`),
--   UNIQUE INDEX `user_mail_UNIQUE` (`user_mail` ASC),
--   UNIQUE INDEX `user_name_UNIQUE` (`user_name` ASC),
--   UNIQUE INDEX `user_id_UNIQUE` (`user_id` ASC))ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- INSERT INTO golang_user ( user_name, user_password, user_registerDate, user_lastLoginTime)
--                        VALUES
--                        ( 'abc', '1', NOW(), NOW() );
