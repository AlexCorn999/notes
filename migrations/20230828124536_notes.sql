-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    notes (
        id SERIAL PRIMARY KEY,
        title varchar(255) NOT NULL,
        content TEXT,
        user_id integer REFERENCES users (id)
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE notes;

-- +goose StatementEnd