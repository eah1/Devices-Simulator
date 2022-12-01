-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         UUID,
    first_name TEXT,
    last_name  TEXT,
    email      TEXT UNIQUE,
    password   TEXT,
    language   TEXT,
    company    TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
