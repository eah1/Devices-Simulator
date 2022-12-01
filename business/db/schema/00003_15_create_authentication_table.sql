-- +goose Up
-- +goose StatementBegin
CREATE TABLE authentications
(
    id         UUID,
    token      TEXT,
    user_id    UUID      NOT NULL,
    valid      BOOLEAN            DEFAULT TRUE,
    login_at   TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    logout_at  TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    CONSTRAINT fk_users
        FOREIGN KEY (user_id)
            REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT unique_authentications_user_id_token UNIQUE (user_id, token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE authentications;
-- +goose StatementEnd
