package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"todoapps/internal/model"
)

func TestPostgresRepo_Create(t *testing.T) {
	tests := []struct {
		name    string
		task    model.Task
		dbErr   error
		wantErr bool
	}{
		{
			name:    "success",
			task:    model.Task{Task: "Buy milk", Done: false},
			dbErr:   nil,
			wantErr: false,
		},
		{
			name:    "db error",
			task:    model.Task{Task: "Fail", Done: false},
			dbErr:   errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer func() { _ = db.Close() }()

			repo := &PostgresRepo{db: db}

			if tc.dbErr != nil {
				mock.ExpectExec("INSERT INTO tasks").
					WithArgs(tc.task.Task, tc.task.Done).
					WillReturnError(tc.dbErr)
			} else {
				mock.ExpectExec("INSERT INTO tasks").
					WithArgs(tc.task.Task, tc.task.Done).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = repo.Create(context.Background(), tc.task)
			if (err != nil) != tc.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresRepo_GetAll(t *testing.T) {
	tests := []struct {
		name      string
		rows      []model.Task
		dbErr     error
		wantErr   bool
		wantCount int
	}{
		{
			name: "multiple rows",
			rows: []model.Task{
				{ID: 1, Task: "Task A", Done: false},
				{ID: 2, Task: "Task B", Done: true},
			},
			wantCount: 2,
		},
		{
			name:      "empty",
			rows:      nil,
			wantCount: 0,
		},
		{
			name:    "query error",
			dbErr:   errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer func() { _ = db.Close() }()

			repo := &PostgresRepo{db: db}

			if tc.dbErr != nil {
				mock.ExpectQuery("SELECT id, task, done FROM tasks").
					WillReturnError(tc.dbErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "task", "done"})
				for _, r := range tc.rows {
					rows.AddRow(r.ID, r.Task, r.Done)
				}
				mock.ExpectQuery("SELECT id, task, done FROM tasks").
					WillReturnRows(rows)
			}

			tasks, err := repo.GetAll(context.Background())
			if (err != nil) != tc.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tc.wantErr)
			}
			if len(tasks) != tc.wantCount {
				t.Errorf(
					"GetAll() returned %d tasks, want %d",
					len(tasks),
					tc.wantCount,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresRepo_GetByID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		row      *model.Task
		wantErr  bool
		wantTask string
	}{
		{
			name:     "found",
			id:       "1",
			row:      &model.Task{ID: 1, Task: "Found", Done: true},
			wantTask: "Found",
		},
		{
			name:    "not found",
			id:      "99",
			row:     nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer func() { _ = db.Close() }()

			repo := &PostgresRepo{db: db}

			if tc.row != nil {
				mock.ExpectQuery("SELECT id, task, done FROM tasks WHERE id = \\$1").
					WithArgs(tc.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "task", "done"}).
						AddRow(tc.row.ID, tc.row.Task, tc.row.Done))
			} else {
				mock.ExpectQuery("SELECT id, task, done FROM tasks WHERE id = \\$1").
					WithArgs(tc.id).
					WillReturnError(sql.ErrNoRows)
			}

			task, err := repo.GetByID(context.Background(), tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tc.wantErr)
			}
			if task.Task != tc.wantTask {
				t.Errorf(
					"GetByID() returned task %s, want %s",
					task.Task,
					tc.wantTask,
				)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresRepo_Update(t *testing.T) {
	tests := []struct {
		name    string
		task    model.Task
		dbErr   error
		wantErr bool
	}{
		{
			name: "success",
			task: model.Task{ID: 1, Task: "Updated", Done: true},
		},
		{
			name:    "db error",
			task:    model.Task{ID: 1, Task: "Fail", Done: false},
			dbErr:   errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer func() { _ = db.Close() }()

			repo := &PostgresRepo{db: db}

			if tc.dbErr != nil {
				mock.ExpectExec("UPDATE tasks SET task = \\$1, done = \\$2 WHERE id = \\$3").
					WithArgs(tc.task.Task, tc.task.Done, tc.task.ID).
					WillReturnError(tc.dbErr)
			} else {
				mock.ExpectExec("UPDATE tasks SET task = \\$1, done = \\$2 WHERE id = \\$3").
					WithArgs(tc.task.Task, tc.task.Done, tc.task.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = repo.Update(context.Background(), tc.task)
			if (err != nil) != tc.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestPostgresRepo_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		dbErr   error
		wantErr bool
	}{
		{
			name: "success",
			id:   "1",
		},
		{
			name:    "db error",
			id:      "5",
			dbErr:   errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer func() { _ = db.Close() }()

			repo := &PostgresRepo{db: db}

			if tc.dbErr != nil {
				mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
					WithArgs(tc.id).
					WillReturnError(tc.dbErr)
			} else {
				mock.ExpectExec("DELETE FROM tasks WHERE id = \\$1").
					WithArgs(tc.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err = repo.Delete(context.Background(), tc.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tc.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
