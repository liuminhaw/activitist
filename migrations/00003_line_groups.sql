-- +goose Up
-- +goose StatementBegin
CREATE TABLE line_groups (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    identifier VARCHAR(40) NOT NULL,
    group_id INT UNIQUE NOT NULL,
    CONSTRAINT fk_groups 
        FOREIGN KEY (group_id) REFERENCES groups(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE line_groups;
-- +goose StatementEnd
