-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE USERS (
    id UUID PRIMARY KEY default gen_random_uuid(),
    name varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE USERS;
-- +goose StatementEnd
