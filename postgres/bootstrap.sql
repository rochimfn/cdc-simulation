CREATE DATABASE cdc_target;

\c cdc_target;

CREATE TABLE person (
    person_id SERIAL PRIMARY KEY,
    fullname varchar(255) NOT NULL,
    num_transaction int,
    company_name varchar(255)
);

INSERT INTO person (person_id, fullname, num_transaction, company_name) VALUES (1, 'John Doe', 1, 'Ajaib');