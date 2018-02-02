--
CREATE DATABASE IF NOT EXISTS `user`;
USE `user`;

-- ----------------------------------------------------------

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(12) unsigned NOT NULL AUTO_INCREMENT,
  `UnionID` text,
  `UserName` varchar(16) NOT NULL DEFAULT '',
  `Password` varchar(128) NOT NULL DEFAULT '',
  `AvatarID` varchar(128) NOT NULL DEFAULT '',
  `Phone` varchar(12) NOT NULL DEFAULT '',
  `IsActive` BOOL NOT NULL DEFAULT TRUE ,
  `ArticleNum` int(64)  NOT NULL DEFAULT 0,
  `Type` INT(11)  NOT NULL DEFAULT 0,
  `Created` datetime NOT NULL DEFAULT current_timestamp,
  `LastLogin` datetime NOT NULL DEFAULT current_timestamp,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
