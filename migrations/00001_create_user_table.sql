-- +goose Up
CREATE TABLE "user" (
    id SERIAL NOT NULL PRIMARY KEY,
    login varchar(255) NOT NULL UNIQUE,
    pass varchar(255) NOT NULL,
    balance NUMERIC(12, 2) NOT NULL DEFAULT 0,
    withdrawn NUMERIC(12, 2) NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE "user";
