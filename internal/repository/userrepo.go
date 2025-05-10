package repository

import (
	"context"
	"errors"

	"github.com/Alheor/gophermart/internal/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	Postgres
}

type UserRepo interface {
	CreateUser(ctx context.Context, form *models.RegistrationForm) (*models.User, error)
	FindUser(ctx context.Context, form *models.LoginForm) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

func GetUserRepository() UserRepo {
	ur := new(UserRepository)
	ur.Conn = connection.Conn
	return ur
}

func (ur *UserRepository) CreateUser(ctx context.Context, form *models.RegistrationForm) (*models.User, error) {

	passHash, err := getPasswordHash(form.Password)
	if err != nil {
		return nil, err
	}

	id := 0
	err = ur.Conn.QueryRow(ctx,
		`INSERT INTO "user" (login, pass) VALUES (@login, @pass) RETURNING id`,
		pgx.NamedArgs{"login": form.Login, "pass": passHash},
	).Scan(&id)

	if err != nil {
		var myErr *pgconn.PgError
		if errors.As(err, &myErr) && myErr.Code == pgerrcode.UniqueViolation {
			return nil, &UniqueErr{}
		}

		return nil, err
	}

	return &models.User{ID: id}, nil
}

func (ur *UserRepository) FindUser(ctx context.Context, form *models.LoginForm) (*models.User, error) {

	var row pgx.Row
	var password string
	var user models.User

	row = ur.Conn.QueryRow(ctx,
		`SELECT id, login, pass FROM "user" WHERE login=@login`,
		pgx.NamedArgs{"login": form.Login},
	)

	err := row.Scan(&user.ID, &user.Login, &password)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(form.Password))
	if err != nil {
		return nil, nil
	}

	return &user, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var row pgx.Row
	var user models.User

	row = ur.Conn.QueryRow(ctx,
		`SELECT id, login, balance, withdrawn FROM "user" WHERE id=@id`,
		pgx.NamedArgs{"id": id},
	)

	err := row.Scan(&user.ID, &user.Login, &user.Balance, &user.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getPasswordHash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}
