-- +goose Up

-- User table
CREATE TABLE "user" (
    id                          INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_full_name              VARCHAR(255) NOT NULL,
    email                       VARCHAR(255) NOT NULL UNIQUE,
    password_hash               VARCHAR(255) NOT NULL,
    phone                       VARCHAR(20),
    requires_password_update    BOOLEAN NOT NULL DEFAULT TRUE,
    role                        VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status                      VARCHAR(20) NOT NULL DEFAULT 'enabled',
    change_history              JSONB,
    stored_in_change_log_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_user_role    CHECK (role IN ('superadmin', 'admin', 'user')),
    CONSTRAINT chk_user_status  CHECK (status IN ('enabled', 'disabled'))
);

-- +goose Down

DROP TABLE "user";
