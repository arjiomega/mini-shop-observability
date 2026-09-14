-- +goose Up
CREATE TABLE products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);


-- +goose Down
DROP TABLE products;
