package repository

import (
	"context"
	"database/sql"

	"todoapps/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, t model.Task) error
	GetAll(ctx context.Context) ([]model.Task, error)
}

type PostgresRepo struct {
	db *sql.DB
}

func (r *PostgresRepo) Create(ctx context.Context, t model.Task) error {
	query := `INSERT INTO tasks (task, done) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, t.Task, t.Done)
	return err
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]model.Task, error) {
	query := `SELECT id, task, done FROM tasks`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Task, &t.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
