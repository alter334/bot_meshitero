CREATE DATABASE IF NOT EXISTS `meshitero`;
USE `meshitero`;
CREATE USER meshitero IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON meshitero.* TO 'meshitero'@'%';
CREATE TABLE IF NOT EXISTS `users` (
  `name` text NOT NULL,
  `id` char(36) NOT NULL,
  `attack` int(11) NOT NULL,
  `rate` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `places` (
  `channelid` char(36) NOT NULL,
  `channelusername` text NOT NULL,
  PRIMARY KEY (`channelid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
