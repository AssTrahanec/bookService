-- +goose Up
CREATE TABLE books
(
    id               UUID PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    author           VARCHAR(255) NOT NULL,
    publication_year INT,
    genre            VARCHAR(100)
);
CREATE TABLE users
       (
       user_id UUID PRIMARY KEY,
       role VARCHAR(255) NOT NULL
       );

-- +goose Down
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS users;
