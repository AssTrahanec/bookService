-- +goose Up
CREATE TABLE books
(
    book_id          UUID PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    author           VARCHAR(255) NOT NULL,
    publication_year SMALLINT,
    genre            VARCHAR(100)
);

CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_books_publication_year ON books(publication_year);
CREATE INDEX idx_books_genre ON books(genre);

CREATE TABLE users
(
    user_id       UUID PRIMARY KEY,
    role          VARCHAR(255) NOT NULL,
    login         VARCHAR(255) NOT NULL,
    password_hash varchar(255) NOT NULL
);
CREATE TABLE users_books
(
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, book_id)
);

CREATE INDEX idx_users_books_book_id ON users_books(book_id);

-- +goose Down
DROP INDEX IF EXISTS idx_books_author;
DROP INDEX IF EXISTS idx_books_publication_year;
DROP INDEX IF EXISTS idx_books_genre;
DROP INDEX IF EXISTS idx_users_books_book_id;
DROP TABLE IF EXISTS users_books;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS users;

