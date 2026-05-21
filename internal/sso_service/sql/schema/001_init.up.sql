-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS person (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL CHECK (octet_length(password_hash) > 1),
    user_role TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_person_role ON person (user_role);

CREATE TABLE IF NOT EXISTS refresh_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES person (id) ON DELETE CASCADE,
    app_id INT NOT NULL,
    token_id BYTEA NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_id ON refresh_sessions (token_id);

INSERT INTO person (
    email,
    username,
    password_hash,
    user_role
) VALUES (
    'eliaseromanov@gmal.com',
    'admin',
    crypt('admin', gen_salt('bf'))::BYTEA,
    'admin'
) ON CONFLICT DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS refresh_sessions;

DROP TABLE IF EXISTS person;
