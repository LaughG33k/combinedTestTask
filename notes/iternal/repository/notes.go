package repository

import (
	"context"

	"github.com/LaughG33k/notes/iternal/model"
	"github.com/jackc/pgx"
)

type Note struct {
	Conn *pgx.ConnPool
}

func (r *Note) Create(ctx context.Context, note model.Note) error {

	if _, err := r.Conn.ExecEx(ctx, "insert into notes(owner_uuid, title, content) values($1, $2, $3);", nil, note.OwnerUuid, note.Title, note.Content); err != nil {
		return err
	}

	return nil
}

func (r *Note) Get(ctx context.Context, owner_uuid string) ([]model.Note, error) {

	res := make([]model.Note, 0)

	rows, err := r.Conn.QueryEx(ctx, "select title, content from notes where owner_uuid = $1", nil, owner_uuid)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var note model.Note

		if err := rows.Scan(&note.Title, &note.Content); err != nil {
			continue
		}

		res = append(res, note)

	}

	return res, nil

}
