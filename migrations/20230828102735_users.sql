-- +goose Up

-- +goose StatementBegin

CREATE TABLE
    users (
        user_id bigserial not null primary key,
        email varchar not null unique,
        encrypted_password varchar not null
    );

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE users;

-- +goose StatementEnd

-- goose postgres "host=localhost user=kode dbname=kode password=5427 sslmode=disable" up