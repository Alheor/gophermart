-- +goose Up
CREATE TABLE "withdrawal" (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_number varchar(255) NOT NULL,
    amount NUMERIC(12, 2) NOT NULL
);

CREATE INDEX withdrawal_user_id_idx ON "withdrawal" (user_id);
CREATE INDEX withdrawal_created_at_idx ON "withdrawal" (created_at);
CREATE UNIQUE INDEX withdrawal_user_id_order_number_key ON "withdrawal" (user_id, order_number);
CREATE UNIQUE INDEX withdrawal_order_number_key ON "withdrawal" (order_number);

ALTER TABLE "withdrawal" ADD CONSTRAINT withdrawal_order_number_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) NOT DEFERRABLE INITIALLY IMMEDIATE;

-- +goose Down
DROP INDEX withdrawal_user_id_order_number_key;
DROP INDEX withdrawal_order_number_key;

DROP TABLE "order";