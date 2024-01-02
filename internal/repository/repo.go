package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alheor/gophermart/internal/config"
)

const (
	uniqueIndexNameByOrderAndUser = `order_user_id_order_number_key`
	uniqueIndexNameByOrder        = `order_order_number_key`
)

// Postgres connection structure
type Postgres struct {
	Conn *pgxpool.Pool
}

type UniqueErr struct {
	error
}

var connection *Postgres

func (pg *Postgres) Init(ctx context.Context) error {

	if pg.Conn != nil {
		return nil
	}

	db, err := pgxpool.New(ctx, config.Options.DatabaseURI)
	if err != nil {
		panic(err)
	}

	pg.Conn = db

	createDBSchema(ctx, pg.Conn)

	return nil
}

func Init(ctx context.Context) error {

	connection = new(Postgres)

	err := connection.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetConnection() *Postgres {
	return connection
}

func createDBSchema(ctx context.Context, conn *pgxpool.Pool) {

	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS "user" (
		    id SERIAL NOT NULL PRIMARY KEY,
			login varchar(255) NOT NULL UNIQUE,
			pass varchar(255) NOT NULL,
			balance NUMERIC(12, 2) NOT NULL DEFAULT 0,
			withdrawn NUMERIC(12, 2) NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS "order" (
		    id SERIAL NOT NULL PRIMARY KEY,
		    user_id INT NOT NULL,
		    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			order_number varchar(255) NOT NULL,
			accrual NUMERIC(12, 2) NOT NULL DEFAULT 0,
			status varchar(10) NOT NULL DEFAULT 'NEW'
		);

	    CREATE INDEX IF NOT EXISTS order_user_id_idx ON "order" (user_id);
		CREATE INDEX IF NOT EXISTS order_created_at_idx ON "order" (created_at);
		CREATE INDEX IF NOT EXISTS order_order_number_idx ON "order" (order_number);
		CREATE UNIQUE INDEX IF NOT EXISTS `+uniqueIndexNameByOrderAndUser+` ON "order" (user_id, order_number);
		CREATE UNIQUE INDEX IF NOT EXISTS `+uniqueIndexNameByOrder+` ON "order" (order_number);

		ALTER TABLE "order" DROP CONSTRAINT IF EXISTS order_order_number_fkey;
		ALTER TABLE "order" ADD CONSTRAINT order_order_number_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) NOT DEFERRABLE INITIALLY IMMEDIATE;

		CREATE TABLE IF NOT EXISTS "withdrawal" (
		    id SERIAL NOT NULL PRIMARY KEY,
		    user_id INT NOT NULL,
		    created_at TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
			order_number varchar(255) NOT NULL,
			amount NUMERIC(12, 2) NOT NULL
		);

 		CREATE INDEX IF NOT EXISTS withdrawal_user_id_idx ON "withdrawal" (user_id);
		CREATE INDEX IF NOT EXISTS withdrawal_created_at_idx ON "withdrawal" (created_at);
		CREATE UNIQUE INDEX IF NOT EXISTS `+uniqueIndexNameByOrderAndUser+` ON "withdrawal" (user_id, order_number);
		CREATE UNIQUE INDEX IF NOT EXISTS `+uniqueIndexNameByOrder+` ON "withdrawal" (order_number);

		ALTER TABLE "withdrawal" DROP CONSTRAINT IF EXISTS withdrawal_order_number_fkey;
		ALTER TABLE "withdrawal" ADD CONSTRAINT withdrawal_order_number_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) NOT DEFERRABLE INITIALLY IMMEDIATE;

	`)

	if err != nil {
		panic(err)
	}
}
