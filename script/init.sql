CREATE DATABASE IF NOT EXISTS `meshitero`;
USE `meshitero`;
CREATE USER meshitero IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON meshitero.* TO 'meshitero'@'%';
CREATE TABLE IF NOT EXISTS `user` (
  `name` text NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `attack` int(11) NOT NULL,
  `rate` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

