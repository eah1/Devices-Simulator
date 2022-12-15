-- +goose Up
-- +goose StatementBegin
CREATE TABLE devices
(
    id               UUID,
    name             TEXT,
    user_id          UUID      NOT NULL,
    environment_id   UUID      NOT NULL,
    device_config_id UUID      NOT NULL,

    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (id),
    CONSTRAINT fk_users
        FOREIGN KEY (user_id)
            REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_environments
        FOREIGN KEY (environment_id)
            REFERENCES environments (id) ON DELETE CASCADE,
    CONSTRAINT fk_devices_config
        FOREIGN KEY (device_config_id)
            REFERENCES devices_config (id) ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE devices;
-- +goose StatementEnd
