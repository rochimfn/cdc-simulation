CREATE DATABASE IF NOT EXISTS `cdc_reference`;

USE `cdc_reference`;

CREATE TABLE IF NOT EXISTS `company` (
    company_id int NOT NULL AUTO_INCREMENT,
    fullname varchar(255) NOT NULL,
    PRIMARY KEY (company_id)
);

INSERT INTO `company` (fullname) VALUES 
("Ajaib"),
("Bubar"),
("Cermat"),
("Dunia");

