-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE USERS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    birthdate DATE,
    bio TEXT,
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now()
);

CREATE TABLE POSTS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    text TEXT,
    user_id UUID,
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE POSTS;
DROP TABLE USERS;
-- +goose StatementEnd
