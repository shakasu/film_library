CREATE TABLE actors
(
    id            serial       not null unique,
    name          varchar(255) not null,
    gender        varchar(1)   not null,
    date_of_birth date         not null
);

CREATE TYPE enum_rating AS enum (
    1,
    2,
    3,
    4,
    5,
    6,
    7,
    8,
    9,
    10
);

CREATE TABLE films
(
    id           serial        not null unique,
    name         varchar(150)  not null,
    description  varchar(1000) not null,
    release_date date          not null,
    rating       enum_rating       not null
);
