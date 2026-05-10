package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"todoapps/internal/model"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, t model.Task) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, t model.Task) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetByID(
	ctx context.Context,
	id string,
) (model.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Task), args.Error(1)
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockTask   model.Task
		mockErr    error
		wantStatus int
	}{
		{
			name:       "success",
			body:       `{"task":"Write tests","done":false}`,
			mockTask:   model.Task{Task: "Write tests", Done: false},
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid json",
			body:       "bad json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "repo error",
			body:       `{"task":"Fail task","done":false}`,
			mockTask:   model.Task{Task: "Fail task", Done: false},
			mockErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &TaskHandler{repo: mockRepo}

			if tc.wantStatus != http.StatusBadRequest {
				mockRepo.On("Create", mock.Anything, tc.mockTask).
					Return(tc.mockErr)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader([]byte(tc.body)),
			)
			rec := httptest.NewRecorder()

			handler.CreateTask(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name       string
		mockTasks  []model.Task
		mockErr    error
		wantStatus int
		wantCount  int
	}{
		{
			name: "success",
			mockTasks: []model.Task{
				{ID: 1, Task: "Task one", Done: false},
				{ID: 2, Task: "Task two", Done: true},
			},
			mockErr:    nil,
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "repo error",
			mockTasks:  []model.Task(nil),
			mockErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &TaskHandler{repo: mockRepo}

			mockRepo.On("GetAll", mock.Anything).
				Return(tc.mockTasks, tc.mockErr)

			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			rec := httptest.NewRecorder()

			handler.GetAllTasks(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantStatus == http.StatusOK {
				assert.Equal(
					t,
					"application/json",
					rec.Header().Get("Content-Type"),
				)
				var result []model.Task
				err := json.NewDecoder(rec.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Len(t, result, tc.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockTask   model.Task
		mockErr    error
		wantStatus int
	}{
		{
			name:       "found",
			id:         "1",
			mockTask:   model.Task{ID: 1, Task: "Read Go docs", Done: false},
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "not found",
			id:         "99",
			mockTask:   model.Task{},
			mockErr:    assert.AnError,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &TaskHandler{repo: mockRepo}

			mockRepo.On("GetByID", mock.Anything, tc.id).
				Return(tc.mockTask, tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /tasks/{id}", handler.GetTask)

			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tc.id, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantStatus == http.StatusOK {
				var result model.Task
				err := json.NewDecoder(rec.Body).Decode(&result)
				assert.NoError(t, err)
				assert.Equal(t, tc.mockTask, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		body       string
		mockTask   model.Task
		mockErr    error
		setupMock  bool
		wantStatus int
	}{
		{
			name:       "success",
			id:         "1",
			body:       `{"task":"Updated task","done":true}`,
			mockTask:   model.Task{ID: 1, Task: "Updated task", Done: true},
			mockErr:    nil,
			setupMock:  true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			id:         "abc",
			body:       `{"task":"task","done":false}`,
			setupMock:  false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			id:         "1",
			body:       "bad",
			setupMock:  false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &TaskHandler{repo: mockRepo}

			if tc.setupMock {
				mockRepo.On("Update", mock.Anything, tc.mockTask).
					Return(tc.mockErr)
			}

			mux := http.NewServeMux()
			mux.HandleFunc("PUT /tasks/{id}", handler.UpdateTask)

			req := httptest.NewRequest(
				http.MethodPut,
				"/tasks/"+tc.id,
				bytes.NewReader([]byte(tc.body)),
			)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mockErr    error
		wantStatus int
	}{
		{
			name:       "success",
			id:         "1",
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "repo error",
			id:         "5",
			mockErr:    assert.AnError,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			handler := &TaskHandler{repo: mockRepo}

			mockRepo.On("Delete", mock.Anything, tc.id).Return(tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /tasks/{id}", handler.DeleteTask)

			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tc.id, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
