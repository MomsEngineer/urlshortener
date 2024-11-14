-- +migrate Up
ALTER TABLE links
ADD COLUMN user_id VARCHAR(255) NOT NULL;