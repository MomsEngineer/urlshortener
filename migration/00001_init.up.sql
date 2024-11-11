-- +migrate Up
CREATE TABLE IF NOT EXISTS links (
    id SERIAL PRIMARY KEY,
    short_link VARCHAR(255) NOT NULL,
    original_link TEXT NOT NULL UNIQUE
);