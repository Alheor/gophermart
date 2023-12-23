package repository

import (
	"context"
	"errors"
	"github.com/Alheor/gophermart/internal/entity"
	"github.com/Alheor/gophermart/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo interface {
	AddOrder(ctx context.Context, user *entity.User, orderNumber string) error
	GetOrders(ctx context.Context, user *entity.User) ([]entity.Order, error)
	GetOrderForProcessing(ctx context.Context) ([]entity.Order, error)
	ChangeOrder(ctx context.Context, order *entity.AccrualOrder) error
}

type OrderRepository struct {
	Conn *pgxpool.Pool
}

func GetOrderRepository() OrderRepo {
	or := new(OrderRepository)
	or.Conn = connection.Conn
	return or
}

func (or *OrderRepository) AddOrder(ctx context.Context, user *entity.User, orderNumber string) error {
	_, err := or.Conn.Exec(ctx,
		`INSERT INTO "order" (user_id, order_number) VALUES (@user_id, @order_number)`,
		pgx.NamedArgs{"user_id": user.ID, "order_number": orderNumber},
	)

	if err == nil {
		return nil
	}

	var myErr *pgconn.PgError
	if errors.As(err, &myErr) && myErr.Code == pgerrcode.UniqueViolation {
		if myErr.ConstraintName == uniqueIndexNameByOrderAndUser {

			return &entity.UniqueErrByUserAndOrder{}
		}

		if myErr.ConstraintName == uniqueIndexNameByOrder {
			return &entity.UniqueErrByOrder{}
		}
	}

	return err
}

func (or *OrderRepository) GetOrders(ctx context.Context, user *entity.User) ([]entity.Order, error) {

	rows, err := or.Conn.Query(ctx,
		`SELECT created_at, order_number, status, accrual FROM "order" WHERE user_id=@id`,
		pgx.NamedArgs{"id": user.ID},
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orderList := make([]entity.Order, 0)

	for rows.Next() {
		var order entity.Order

		err = rows.Scan(&order.UploadedAt, &order.Number, &order.Status, &order.Accrual)
		if err != nil {
			return nil, err
		}

		orderList = append(orderList, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orderList, nil
}

func (or *OrderRepository) GetOrderForProcessing(ctx context.Context) ([]entity.Order, error) {
	rows, err := or.Conn.Query(ctx,
		`SELECT order_number, status FROM "order" WHERE status NOT IN(@status_invalid, @status_processed)`,
		pgx.NamedArgs{"status_invalid": "INVALID", "status_processed": "PROCESSED"},
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orderList := make([]entity.Order, 0)
	for rows.Next() {
		var order entity.Order

		err = rows.Scan(&order.Number, &order.Status)
		if err != nil {
			return nil, err
		}

		orderList = append(orderList, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orderList, nil
}

func (or *OrderRepository) ChangeOrder(ctx context.Context, order *entity.AccrualOrder) error {

	tx, err := or.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.GetLogger().Error(`Change order error1: ` + err.Error())
		return err
	}
	logger.GetLogger().Error(`Change order error2: `)
	if order.Accrual > 0 {
		_, err = or.Conn.Exec(ctx,
			`UPDATE "order" SET status = @status, accrual = @accrual WHERE order_number = @order_number;`,
			pgx.NamedArgs{"status": order.Status, "accrual": order.Accrual, "order_number": order.Order},
		)

		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return err
			}

			return err
		}

		_, err = or.Conn.Exec(ctx,
			`
					UPDATE "user" SET balance = balance + @accrual WHERE id = (
						SELECT user_id FROM "order" WHERE order_number = @order_number LIMIT 1
					);
				`,
			pgx.NamedArgs{"status": order.Status, "accrual": order.Accrual, "order_number": order.Order},
		)

		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return err
			}

			return err
		}

	} else {
		_, err = or.Conn.Exec(ctx,
			`UPDATE "order" SET status = @status WHERE order_number = @order_number;`,
			pgx.NamedArgs{"status": order.Status, "order_number": order.Order},
		)

		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				return err
			}

			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
