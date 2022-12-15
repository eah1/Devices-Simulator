-- +goose Up
-- +goose StatementBegin
CREATE TABLE devices_config
(
    id                  UUID,
    name                TEXT,
    vars                JSONB              DEFAULT '{}',
    metrics_fixed       JSONB              DEFAULT '{}',
    metrics_accumulated JSONB              DEFAULT '{}',
    type_send           TEXT,
    payload             TEXT,
    user_id             UUID      NOT NULL,

    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    CONSTRAINT fk_users
        FOREIGN KEY (user_id)
            REFERENCES users (id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE devices_config;
-- +goose StatementEnd
