package repository

import (
	"context"

	"github.com/Alheor/gophermart/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	uniqueIndexNameByOrderAndUser = `order_user_id_order_number_key`
	uniqueIndexNameByOrder        = `order_order_number_key`
)

type UniqueErr struct {
	error
}

// Postgres connection structure
type Postgres struct {
	Conn *pgxpool.Pool
}

var connection *Postgres

func Init(ctx context.Context, databaseURI string) error {

	connection = new(Postgres)

	if connection.Conn != nil {
		return nil
	}

	db, err := pgxpool.New(ctx, databaseURI)
	if err != nil {
		return err
	}

	logger.Info(`Running migrations ...`)

	//./../../migrations - нужно, что бы работали тесты
	if err = goose.Up(stdlib.OpenDBFromPool(db), "./migrations"); err != nil {
		if err = goose.Up(stdlib.OpenDBFromPool(db), "./../../migrations"); err != nil {
			logger.Error(`run migrations error: `, err)
		}
	}

	connection.Conn = db

	return nil
}

func GetConnection() *Postgres {
	return connection
}
