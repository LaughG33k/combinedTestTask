package psql

import (
	"context"
	"time"

	"github.com/LaughG33k/authServiceTestTask/iternal"
	"github.com/LaughG33k/authServiceTestTask/pkg"
	"github.com/jackc/pgx"
)

func NewClient(ctx context.Context, cfg iternal.DBConfig) (*pgx.ConnPool, error) {

	var pool *pgx.ConnPool

	err := pkg.Rerty(func() error {

		p, err := pgx.NewConnPool(pgx.ConnPoolConfig{
			MaxConnections: cfg.PoolSize,
			ConnConfig: pgx.ConnConfig{
				User:     cfg.User,
				Password: cfg.Password,
				Host:     cfg.Host,
				Port:     cfg.Port,
				Database: cfg.DB,
			},
		})

		if err != nil {
			return err
		}

		pool = p

		return nil

	}, cfg.TryAttempts, 3*time.Second)

	if err != nil {

		return nil, err

	}

	return pool, nil

}
