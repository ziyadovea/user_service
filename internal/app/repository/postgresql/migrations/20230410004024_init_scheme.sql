-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(100),
    email VARCHAR(100),
    password VARCHAR(100),

    PRIMARY KEY (id),
    UNIQUE (name),
    UNIQUE (email)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
