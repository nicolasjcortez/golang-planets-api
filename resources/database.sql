CREATE DATABASE starwars;
USE starwars;

DROP TABLE IF EXISTS `planets`;
CREATE TABLE `planets` (
  `planet_id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `terrain` varchar(100) NOT NULL,
  `climate` varchar(100) NOT NULL,
  `qtd_films` int NOT NULL,
  PRIMARY KEY (`planet_id`)
)
