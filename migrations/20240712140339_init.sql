-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE USERS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    date_of_birth DATE,
    bio TEXT,
    avatar_url VARCHAR(255),
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now()
);

CREATE TABLE POSTS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    photo_url VARCHAR(255),
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE
);

CREATE TABLE POST_LIKES(
    id uuid PRIMARY KEY default gen_random_uuid(),
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE,
    FOREIGN KEY("post_id") REFERENCES POSTS("id") ON DELETE CASCADE
);

CREATE TABLE COMMENTS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    parent_comment_id UUID,
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE,
    FOREIGN KEY("post_id") REFERENCES POSTS("id") ON DELETE CASCADE,
    FOREIGN KEY("parent_comment_id") REFERENCES COMMENTS("id") ON DELETE CASCADE
);

CREATE TABLE COMMENT_LIKES (
    id uuid PRIMARY KEY default gen_random_uuid(),
    comment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    updated_at TIMESTAMP NOT NULL default now(),
    created_at TIMESTAMP NOT NULL default now(),
    FOREIGN KEY("user_id") REFERENCES USERS("id") ON DELETE CASCADE,
    FOREIGN KEY("comment_id") REFERENCES COMMENTS("id") ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE POST_LIKES;
DROP TABLE COMMENT_LIKES;
DROP TABLE COMMENTS;
DROP TABLE POSTS;
DROP TABLE USERS;
-- +goose StatementEnd
