CREATE ROLE read_only_access;

GRANT CONNECT ON DATABASE postgres
    TO read_only_access;

GRANT USAGE ON SCHEMA public TO read_only_access;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO read_only_access;

CREATE USER read_user WITH
    PASSWORD 'read_user';

GRANT read_only_access TO read_user;

CREATE USER admin WITH PASSWORD 'admin' SUPERUSER;

CREATE TABLE IF NOT EXISTS actors
(
    id            serial PRIMARY KEY,
    name          varchar(255) not null,
    gender        varchar(1)   not null,
    date_of_birth date         not null
);

CREATE TABLE IF NOT EXISTS films
(
    id           serial PRIMARY KEY,
    name         varchar(150)  not null,
    description  varchar(1000) not null,
    release_date date          not null,
    rating       smallserial   not null
);

CREATE TABLE IF NOT EXISTS actor_film
(
    id       SERIAL PRIMARY KEY,
    actor_id INTEGER NOT NULL REFERENCES actors,
    film_id  INTEGER NOT NULL REFERENCES films,
    UNIQUE (actor_id, film_id)
);
