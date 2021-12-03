-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS invite_codes(
    id BIGSERIAL PRIMARY KEY,
    code text UNIQUE NOT NULL,
    valid boolean NOT NULL DEFAULT true,
    created_at timestamp WITH time zone DEFAULT now()
);
CREATE INDEX invite_codes_idx on invite_codes(code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invite_codes;
-- +goose StatementEnd
