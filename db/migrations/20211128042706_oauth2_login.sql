-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS oauth2_login(
    id BIGSERIAL PRIMARY KEY,
    service text NOT NULL,
    email text NOT NULL UNIQUE,
    access_token text NOT NULL UNIQUE,
    refresh_token text NOT NULL UNIQUE,
    created_at timestamp WITH time zone DEFAULT now()
);
CREATE INDEX oauth2_login_idx on oauth2_login(service, email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS oauth2_login;
-- +goose StatementEnd
