package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Alheor/gophermart/internal/models"
)

type WithdrawalOrderRepo interface {
	AddWithdrawOrder(ctx context.Context, user *models.User, form *models.UserWithdrawOrder) error
	GetWithdrawals(ctx context.Context, user *models.User) ([]models.WithdrawalOrder, error)
}

type WithdrawalOrderRepository struct {
	Conn *pgxpool.Pool
}

func GetWithdrawOrderRepository() WithdrawalOrderRepo {
	or := new(WithdrawalOrderRepository)
	or.Conn = connection.Conn
	return or
}

func (or *WithdrawalOrderRepository) AddWithdrawOrder(ctx context.Context, user *models.User, form *models.UserWithdrawOrder) error {

	tx, err := or.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	var row pgx.Row
	var balance float32
	row = or.Conn.QueryRow(ctx,
		`SELECT balance FROM "user" WHERE id = @user_id FOR UPDATE`,
		pgx.NamedArgs{"user_id": user.ID},
	)

	err = row.Scan(&balance)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}

		return err
	}

	if form.Sum > balance {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}

		return &models.ErrNotEnoughMemory{}
	}

	_, err = or.Conn.Exec(ctx,
		`INSERT INTO "withdrawal" (user_id, order_number, amount) VALUES (@user_id, @order_number, @amount)`,
		pgx.NamedArgs{"user_id": user.ID, "order_number": form.Order, "amount": form.Sum},
	)

	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}

		return err
	}

	_, err = or.Conn.Exec(ctx,
		`UPDATE "user" SET withdrawn = withdrawn + @sum, balance = balance - @sum WHERE id = @user_id;`,
		pgx.NamedArgs{"sum": form.Sum, "user_id": user.ID},
	)

	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			return err
		}

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (or *WithdrawalOrderRepository) GetWithdrawals(ctx context.Context, user *models.User) ([]models.WithdrawalOrder, error) {
	rows, err := or.Conn.Query(ctx,
		`SELECT created_at, order_number, amount FROM "withdrawal" WHERE user_id=@id ORDER BY created_at DESC`,
		pgx.NamedArgs{"id": user.ID},
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orderList := make([]models.WithdrawalOrder, 0)

	for rows.Next() {
		var order models.WithdrawalOrder

		err = rows.Scan(&order.ProcessedAt, &order.Order, &order.Sum)
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
