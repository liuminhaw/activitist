-- +goose Up
-- +goose StatementBegin
CREATE TABLE registrations (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    token_hash TEXT NOT NULL UNIQUE,
    expires TIMESTAMPTZ NOT NULL,
    line_id VARCHAR(40) UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE registrations;
-- +goose StatementEnd
