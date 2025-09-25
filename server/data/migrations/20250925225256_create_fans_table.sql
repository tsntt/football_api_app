-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS fans (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, team)
);

CREATE INDEX idx_fans_user_id ON fans(user_id);
CREATE INDEX idx_fans_team ON fans(team);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE fans;
-- +goose StatementEnd
