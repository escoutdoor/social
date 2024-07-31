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
    avatar_url VARCHAR(255),
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now()
);

CREATE TABLE POSTS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    text TEXT NOT NULL,
    user_id UUID NOT NULL,
    photo_url VARCHAR(255),
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE
);

CREATE TABLE REPLIES (
    id UUID PRIMARY KEY default gen_random_uuid(),
    text TEXT NOT NULL,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE,
    FOREIGN KEY("post_id") REFERENCES POSTS("id") ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE REPLIES;
DROP TABLE POSTS;
DROP TABLE USERS;
-- +goose StatementEnd
