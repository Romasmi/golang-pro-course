-- +goose Up
ALTER TABLE events ADD COLUMN notified BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE events DROP COLUMN notified;
