-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts(
    id SERIAL NOT NULL,
    title TEXT NOT NULL UNIQUE,
    author TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT,

    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
