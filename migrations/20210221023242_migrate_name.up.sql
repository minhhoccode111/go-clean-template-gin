CREATE TABLE IF NOT EXISTS history(
    id serial PRIMARY KEY,
    source VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    original VARCHAR(255) NOT NULL,
    translation VARCHAR(255) NOT NULL
);