package repository

import (
	"context"
	"database/sql"

	"todoapps/internal/model"
)

type TaskRepository interface {
	Create(ctx context.Context, t model.Task) error
	GetAll(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id string) (model.Task, error)
	Update(ctx context.Context, t model.Task) error
	Delete(ctx context.Context, id string) error
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresRepo{db: db}, nil
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

func (r *PostgresRepo) GetByID(
	ctx context.Context,
	id string,
) (model.Task, error) {
	var t model.Task
	query := `SELECT id, task, done FROM tasks WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Task, &t.Done)
	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

func (r *PostgresRepo) Update(ctx context.Context, t model.Task) error {
	query := `UPDATE tasks SET task = $1, done = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, t.Task, t.Done, t.ID)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
