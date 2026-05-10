package repository

import (
	"context"
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
