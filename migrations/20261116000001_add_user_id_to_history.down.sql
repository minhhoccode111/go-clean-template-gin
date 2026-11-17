-- DESTRUCTIVE ROLLBACK: drops all per-user translation history and recreates
-- the legacy table without user_id. Data cannot be recovered.
DROP TABLE IF EXISTS history;

CREATE TABLE IF NOT EXISTS history(
    id serial PRIMARY KEY,
    source VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    original VARCHAR(255) NOT NULL,
    translation VARCHAR(255) NOT NULL
);
