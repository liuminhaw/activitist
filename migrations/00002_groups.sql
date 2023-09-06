-- +goose Up
-- +goose StatementBegin
CREATE TABLE groups (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category VARCHAR(32) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE groups;
-- +goose StatementEnd
