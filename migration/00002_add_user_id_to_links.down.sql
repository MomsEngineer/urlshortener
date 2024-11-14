-- +migrate Down
ALTER TABLE links
DROP COLUMN user_id;