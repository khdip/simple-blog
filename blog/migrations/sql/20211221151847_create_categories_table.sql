-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories(
    category_id SERIAL NOT NULL,
    category_name TEXT NOT NULL UNIQUE,

    PRIMARY KEY(category_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
