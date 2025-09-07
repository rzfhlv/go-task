package task_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/repository/task"
	"github.com/rzfhlv/go-task/pkg/param"
	"github.com/stretchr/testify/assert"
)

var (
	now = time.Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC)

	taskModel = model.Task{
		ID:          1,
		Title:       "Todo 1",
		Description: "urgent task",
		Status:      "todo",
		UserID:      int64(1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskModelArr = []model.Task{
		{
			ID:          1,
			Title:       "Todo 1",
			Description: "urgent task",
			Status:      "todo",
			UserID:      int64(1),
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	paramPkg = param.Param{
		Page:   1,
		Limit:  10,
		Offset: 0,
		Total:  0,
	}
)

func TestTaskCreate(t *testing.T) {
	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "user_id", "created_at", "updated_at"}).
					AddRow(taskModel.ID, taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UserID, taskModel.CreatedAt, taskModel.UpdatedAt)

				s.ExpectQuery(`INSERT INTO tasks 
				(title, description, status, user_id)
				VALUES ($1, $2, $3, $4) RETURNING *`).
					WithArgs(taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UserID).
					WillReturnRows(rows)
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when create task",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(`INSERT INTO tasks 
				(title, description, status, user_id)
				VALUES ($1, $2, $3, $4) RETURNING *`).
					WithArgs(taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UserID).
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: model.Task{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := task.New(db)
			result, err := r.Create(context.Background(), taskModel)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestTaskGetByUserID(t *testing.T) {
	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult []model.Task
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "user_id", "created_at", "updated_at"}).
					AddRow(taskModel.ID, taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UserID, taskModel.CreatedAt, taskModel.UpdatedAt)

				s.ExpectQuery(`SELECT 
					id, title, description, status, created_at, updated_at
					FROM tasks
					WHERE user_id = $1
					ORDER BY id LIMIT $2 OFFSET $3`).
					WithArgs(taskModel.UserID, 10, 0).
					WillReturnRows(rows)
			},
			wantResult: taskModelArr,
			wantErr:    nil,
		},
		{
			name: "error when get by id",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(`SELECT 
					id, title, description, status, created_at, updated_at
					FROM tasks
					WHERE user_id = $1
					ORDER BY id LIMIT $2 OFFSET $3`).
					WithArgs(taskModel.UserID, paramPkg.Limit, paramPkg.Offset).
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: []model.Task{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := task.New(db)
			result, err := r.GetByUserID(context.Background(), taskModel.UserID, paramPkg)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestTaskGetByID(t *testing.T) {
	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status", "user_id", "created_at", "updated_at"}).
					AddRow(taskModel.ID, taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UserID, taskModel.CreatedAt, taskModel.UpdatedAt)

				s.ExpectQuery(`SELECT 
					id, title, description, status, created_at, updated_at
					FROM tasks
					WHERE id = $1 AND user_id = $2`).
					WithArgs(taskModel.ID, taskModel.UserID).
					WillReturnRows(rows)
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when get by id",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(`SELECT 
					id, title, description, status, created_at, updated_at
					FROM tasks
					WHERE id = $1 AND user_id = $2`).
					WithArgs(taskModel.ID, taskModel.UserID).
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: model.Task{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := task.New(db)
			result, err := r.GetByID(context.Background(), taskModel.ID, taskModel.UserID)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestTaskUpdate(t *testing.T) {
	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`UPDATE tasks
					SET title = $1, description = $2, status = $3, updated_at = $4
					WHERE id = $5 AND user_id = $6`).
					WithArgs(taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UpdatedAt, taskModel.ID, taskModel.UserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when no rows affected",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`UPDATE tasks
					SET title = $1, description = $2, status = $3, updated_at = $4
					WHERE id = $5 AND user_id = $6`).
					WithArgs(taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UpdatedAt, taskModel.ID, taskModel.UserID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			wantResult: model.Task{},
			wantErr:    errors.New("rows affected error"),
		},
		{
			name: "error when update task",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`UPDATE tasks
					SET title = $1, description = $2, status = $3, updated_at = $4
					WHERE id = $5 AND user_id = $6`).
					WithArgs(taskModel.Title, taskModel.Description, taskModel.Status, taskModel.UpdatedAt, taskModel.ID, taskModel.UserID).
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: model.Task{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := task.New(db)
			result, err := r.Update(context.Background(), taskModel, taskModel.UserID)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestTaskDelete(t *testing.T) {
	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`).
					WithArgs(taskModel.ID, taskModel.UserID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when no rows affected",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`).
					WithArgs(taskModel.ID, taskModel.UserID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			wantResult: model.Task{},
			wantErr:    errors.New("rows affected error"),
		},
		{
			name: "error when update task",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`).
					WithArgs(taskModel.ID, taskModel.UserID).
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: model.Task{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := task.New(db)
			err := r.Delete(context.Background(), taskModel.ID, taskModel.UserID)

			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestTaskCount(t *testing.T) {
	expectedCount := int64(10)
	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult int64
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(expectedCount)

				s.ExpectQuery("SELECT count(*) FROM tasks;").
					WillReturnRows(rows)
			},
			wantResult: 10,
			wantErr:    nil,
		},
		{
			name: "error when count",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT count(*) FROM tasks;").
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: 0,
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := task.New(db)
			result, err := r.Count(context.Background())

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}
