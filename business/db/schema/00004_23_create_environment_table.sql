-- +goose Up
-- +goose StatementBegin
CREATE TABLE environments
(
    id         UUID,
    name       TEXT,
    vars       JSONB              DEFAULT '{}',
    user_id    UUID      NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    CONSTRAINT fk_users
        FOREIGN KEY (user_id)
            REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE environments;
-- +goose StatementEnd
