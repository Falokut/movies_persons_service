CREATE ROLE movies_people_service WITH
    LOGIN
    ENCRYPTED PASSWORD 'SCRAM-SHA-256$4096:R9TMUdvkUG5yxu0rJlO+hA==$E/WRNMfl6SWK9xreXN8rfIkJjpQhWO8pd+8t2kx12D0=:sCS47DCNVIZYhoue/BReTE0ZhVRXzMGszsnnHexVwOU=';

CREATE TABLE people (
    id SERIAL PRIMARY KEY,
    fullname_ru TEXT NOT NULL,
    fullname_en TEXT
);

GRANT SELECT ON people TO movies_people_service;