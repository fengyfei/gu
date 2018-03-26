CREATE TABLE IF NOT EXISTS staff (
  id int(16) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(32) NOT NULL DEFAULT '',
  pwd varchar(128) NOT NULL,
  realname varchar(32) NOT NULL,
  createdat datetime NOT NULL DEFAULT current_timestamp,
  resignat datetime NOT NULL DEFAULT current_timestamp,
  mobile varchar(128),
  email varchar(128),
  male varchar(32),
  active varchar(32),
  resigned varchar(32),
  PRIMARY KEY (id),
  UNIQUE (name)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;