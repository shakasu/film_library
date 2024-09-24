CREATE TABLE IF NOT EXISTS accounts
(
    id        SERIAL,
    login     VARCHAR(10),
    password  VARCHAR(100),
    user_role varchar(6)
);

CREATE EXTENSION pgcrypto;

INSERT INTO accounts (login, password, user_role)
VALUES ('admin',
        crypt('admin', gen_salt('md5')),
        'ADMIN');

INSERT INTO accounts (login, password, user_role)
VALUES ('reader',
        crypt('reader', gen_salt('md5')),
        'READER');

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
