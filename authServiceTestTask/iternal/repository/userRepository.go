package repository

import (
	"context"

	customerrors "github.com/LaughG33k/authServiceTestTask/iternal/errors"
	"github.com/LaughG33k/authServiceTestTask/iternal/model"
	"github.com/jackc/pgx"
)

type UserRepository struct {
	Conn *pgx.ConnPool
}

func (r *UserRepository) CreateUser(ctx context.Context, model model.RegistrationModel) error {

	if _, err := r.Conn.ExecEx(ctx, "insert into users(login, password, name, email) values($1, crypt($2, gen_salt('md5')), $3, $4);", nil, model.Login, model.Password, model.Name, model.Email); err != nil {

		if isErrCode(err, "23505") {
			return customerrors.UserAlreadyExists
		}

		return err

	}

	return nil

}

func (r *UserRepository) CheckLP(ctx context.Context, login, password string) (bool, string, error) {

	exist := false
	uuid := ""

	if err := r.Conn.QueryRowEx(ctx, "select uuid, (password = crypt($2, password)) as t from users where login=$1;", nil, login, password).Scan(&uuid, &exist); err != nil {
		if err == pgx.ErrNoRows {
			return exist, "", customerrors.UserNotFound
		}

		return exist, "", err
	}

	return exist, uuid, nil
}

func (r *UserRepository) GetEmail(ctx context.Context, uuid string) (string, error) {

	var email string

	if err := r.Conn.QueryRowEx(ctx, "select email from refresh_sessions where uuid = $1", nil, uuid).Scan(&email); err != nil {
		if err == pgx.ErrNoRows {
			return "", customerrors.UserNotFound
		}

		return "", err
	}

	return email, nil

}
