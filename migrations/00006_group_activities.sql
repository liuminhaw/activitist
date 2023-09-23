-- +goose Up
-- +goose StatementBegin
CREATE TABLE group_activities (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    activity VARCHAR(255) NOT NULL,
    starttime TIMESTAMPTZ,
    endtime TIMESTAMPTZ,
    location TEXT,
    group_id INT NOT NULL,
    CONSTRAINT fk_groups
        FOREIGN KEY (group_id) REFERENCES groups(id)
        ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE group_activities;
-- +goose StatementEnd
