DROP TABLE IF EXISTS history;

CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    original VARCHAR(255) NOT NULL,
    translation VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_history_user_id ON history(user_id);
