-- +goose Up
-- +goose StatementBegin
CREATE TABLE individual_activities (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    activity VARCHAR(255) NOT NULL,
    starttime TIMESTAMPTZ,
    endtime TIMESTAMPTZ,
    location TEXT,
    user_id INT NOT NULL,
    CONSTRAINT fk_users
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE individual_activities;
-- +goose StatementEnd
