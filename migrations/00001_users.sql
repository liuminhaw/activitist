-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    lineId VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
