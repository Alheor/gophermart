package repository

import (
	"context"
	"errors"
	"github.com/Alheor/gophermart/internal/entity"
	"github.com/Alheor/gophermart/internal/request"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	//ErrUserAlreadyExist error message
	ErrUserAlreadyExist = `User already exist`
)

type UserRepository struct {
	Conn *pgxpool.Pool
}

type UserRepo interface {
	CreateUser(ctx context.Context, form *request.RegisterForm) error
	FindUser(ctx context.Context, form *request.LoginForm) (*entity.User, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
}

func GetUserRepository() UserRepo {
	ur := new(UserRepository)
	ur.Conn = connection.Conn
	return ur
}

func (ur *UserRepository) CreateUser(ctx context.Context, form *request.RegisterForm) error {

	passHash, err := getPasswordHash(form.Password)
	if err != nil {
		return err
	}

	_, err = ur.Conn.Exec(ctx,
		`INSERT INTO "user" (login, pass) VALUES (@login, @pass)`,
		pgx.NamedArgs{"login": form.Login, "pass": passHash},
	)

	if err == nil {
		return nil
	}

	var myErr *pgconn.PgError
	if errors.As(err, &myErr) && myErr.Code == pgerrcode.UniqueViolation {
		return &UniqueErr{}
	}

	return err
}

func (ur *UserRepository) FindUser(ctx context.Context, form *request.LoginForm) (*entity.User, error) {

	var row pgx.Row
	var password string
	var user entity.User

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

func (ur *UserRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {

	var row pgx.Row
	var user entity.User

	row = ur.Conn.QueryRow(ctx,
		`SELECT id, login, balance, withdrawn FROM "user" WHERE id=@id`,
		pgx.NamedArgs{"id": userID},
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
