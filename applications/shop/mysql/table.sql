CREATE DATABASE IF NOT EXISTS `shop`;
USE `shop`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(32) unsigned NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(128) UNIQUE NOT NULL,
  `nickname` VARCHAR(30) NOT NULL,
  `phone` VARCHAR (16) UNIQUE NOT NULL,
  `avatar` VARCHAR(128),
  `password` varchar(128) NOT NULL,
  `sex` INT(8),
  `type` varchar(30) NOT NULL,
  `isadmin` BOOLEAN NOT NULL,
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `address` (
  `id` int(32) unsigned NOT NULL AUTO_INCREMENT,
  `userid` int(32) NOT NULL,
  `address` VARCHAR(128) NOT NULL,
  `phone` VARCHAR (16) NOT NULL,
  `isdefault` BOOLEAN NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;