-- +goose Up
CREATE TABLE "order" (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_number varchar(255) NOT NULL,
    accrual NUMERIC(12, 2) NOT NULL DEFAULT 0,
    status varchar(10) NOT NULL DEFAULT 'NEW'
);

CREATE INDEX order_user_id_idx ON "order" (user_id);
CREATE INDEX order_created_at_idx ON "order" (created_at);
CREATE INDEX order_order_number_idx ON "order" (order_number);
CREATE UNIQUE INDEX order_user_id_order_number_key ON "order" (user_id, order_number);
CREATE UNIQUE INDEX order_order_number_key ON "order" (order_number);

ALTER TABLE "order" ADD CONSTRAINT order_order_number_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) NOT DEFERRABLE INITIALLY IMMEDIATE;

-- +goose Down
DROP INDEX order_user_id_order_number_key;
DROP INDEX order_order_number_key;

DROP TABLE "order";