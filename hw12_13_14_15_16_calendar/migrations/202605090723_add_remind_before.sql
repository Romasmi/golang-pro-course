-- +goose Up
ALTER TABLE events ADD COLUMN remind_before INT NOT NULL DEFAULT 0;
ALTER TABLE events ADD COLUMN notified BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX idx_events_remind ON events (start_at, remind_before)
    WHERE notified = FALSE;

-- +goose Down
ALTER TABLE events DROP COLUMN remind_before;
ALTER TABLE events DROP COLUMN notified;
