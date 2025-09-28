-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS broadcasted_messages (
    id SERIAL PRIMARY KEY,
    match_id INTEGER NOT NULL,
    message_content_hash varchar(64) NOT NULL UNIQUE,
    status VARCHAR(20) DEFAULT 'sent',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_duplicate_check ON broadcasted_messages(match_id, message_content_hash, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE broadcasted_messages;
-- +goose StatementEnd
