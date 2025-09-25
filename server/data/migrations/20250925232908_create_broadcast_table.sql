-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS broadcast_messages (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'sent',
    UNIQUE(match_id)
);

CREATE INDEX idx_broadcast_match_id ON broadcast_messages(match_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE broadcast_messages;
-- +goose StatementEnd
