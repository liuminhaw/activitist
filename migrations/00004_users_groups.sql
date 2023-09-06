-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_groups (
    user_id INT NOT NULL,
    group_id INT NOT NULL,
    PRIMARY KEY (user_id, group_id),
    CONSTRAINT fk_users
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_groups
        FOREIGN KEY (group_id) REFERENCES groups(id)
        ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_groups;
-- +goose StatementEnd
