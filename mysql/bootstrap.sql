CREATE DATABASE IF NOT EXISTS `cdc_source`;

USE `cdc_source`;

CREATE TABLE IF NOT EXISTS `person` (
    person_id int NOT NULL AUTO_INCREMENT,
    fullname varchar(255) NOT NULL,
    num_transaction int,
    company_id int, 
    PRIMARY KEY (person_id)
);

INSERT INTO `person` (person_id, fullname, num_transaction, company_id) VALUES (1, "John Doe", 1, 1);

CREATE USER 'dbz'@'%' IDENTIFIED BY 'password';

GRANT SELECT, RELOAD, SHOW DATABASES, REPLICATION SLAVE, REPLICATION CLIENT, LOCK TABLES ON *.* TO 'dbz'@'%';

FLUSH PRIVILEGES;