CREATE DATABASE IF NOT EXISTS `shop`;
USE `shop`;

CREATE TABLE IF NOT EXISTS `user` (
  `id` INT(64) UNSIGNED NOT NULL AUTO_INCREMENT,
  `openid` VARCHAR(28)  UNIQUE NOT NULL,
  `unionid` VARCHAR(29)  UNIQUE NOT NULL,
  `username` VARCHAR(128),
  `phone` VARCHAR (16),
  `password` VARCHAR(128),
  `avatar` VARCHAR(128),
  `sex` INT(8) UNSIGNED,
  `isadmin` BOOLEAN NOT NULL,
  `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `address` (
  `id` INT(64) UNSIGNED NOT NULL AUTO_INCREMENT,
  `userid` INT(64) UNSIGNED NOT NULL,
  `name` VARCHAR(32) NOT NULL,
  `phone` VARCHAR (16) NOT NULL,
  `address` VARCHAR(128) NOT NULL,
  `isdefault` BOOLEAN NOT NULL DEFAULT FALSE,
  `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE IF NOT EXISTS `category` (
  `id` INT(32) UNSIGNED NOT NULL AUTO_INCREMENT,
  `category` VARCHAR(32) NOT NULL,
  `pid` INT(32) UNSIGNED,
  `status` INT(8) UNSIGNED NOT NULL DEFAULT 0,
  `created` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
