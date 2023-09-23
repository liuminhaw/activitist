-- +goose Up
-- +goose StatementBegin
CREATE TABLE raw_messages (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    message TEXT NOT NULL,
    response JSON NOT NULL,
    user_id INT,
    group_id INT,
    CONSTRAINT fk_users
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE SET NULL,
    CONSTRAINT fk_groups
        FOREIGN KEY (group_id) REFERENCES groups(id)
        ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE raw_messages;
-- +goose StatementEnd
