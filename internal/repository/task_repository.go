package repository

import (
	"context"
	"database/sql"

	"todoapps/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, t model.Task) error
}

type PostgresRepo struct {
	db *sql.DB
}

func (r *PostgresRepo) Create(ctx context.Context, t model.Task) error {
	query := `INSERT INTO tasks (task, done) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, t.Task, t.Done)
	return err
}
