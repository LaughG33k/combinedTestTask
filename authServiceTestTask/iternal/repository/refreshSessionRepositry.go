package repository

import (
	"context"

	customerrors "github.com/LaughG33k/authServiceTestTask/iternal/errors"
	"github.com/jackc/pgx"
)

type RefreshSessionRepository struct {
	Conn *pgx.ConnPool
}

func (r *RefreshSessionRepository) Create(ctx context.Context, ownerUuid, ip, token string, timelife int64) error {

	if _, err := r.Conn.ExecEx(ctx, "delete from refresh_sessions where owner_uuid = $1 and ip = $2;", nil, ownerUuid, ip); err != nil {
		return err
	}

	if _, err := r.Conn.ExecEx(ctx, "insert into refresh_sessions(token, owner_uuid, time_life, ip) values($1, $2, $3, $4)", nil, token, ownerUuid, timelife, ip); err != nil {
		return err
	}

	return nil

}

func (r *RefreshSessionRepository) Update(ctx context.Context, ownerUuid, ip, old, token string, timelife int64) error {

	ctg, err := r.Conn.ExecEx(ctx, "update refresh_sessions set token = $1, time_life = $2 where owner_uuid = $3 and ip = $4 and token = $5;", nil, token, timelife, ownerUuid, ip, old)

	if err != nil {
		return err
	}

	if ctg.RowsAffected() == 0 {
		return customerrors.RefreshNotFound
	}

	return nil
}

func (r *RefreshSessionRepository) Delete(ctx context.Context, ownerUuid, reqtoken string) error {

	if _, err := r.Conn.ExecEx(ctx, "delete from refresh_sessions where owner_uuid = $1 and token = $2;", nil, ownerUuid, reqtoken); err != nil {
		return err
	}

	return nil
}

func (r *RefreshSessionRepository) Get(ctx context.Context, ownerUuid, reqtoken string) (string, int64, error) {

	var ip string
	var timelife int64

	if err := r.Conn.QueryRowEx(ctx, "select ip, time_life from refresh_sessions where owner_uuid = $1 and token = $2;", nil, ownerUuid, reqtoken).Scan(&ip, &timelife); err != nil {

		if err == pgx.ErrNoRows {
			return "", 0, customerrors.RefreshNotFound
		}

		return "", 0, err
	}

	return ip, timelife, nil

}
