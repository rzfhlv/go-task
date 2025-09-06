package task_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/rzfhlv/go-task/internal/model"
	taskmocks "github.com/rzfhlv/go-task/internal/repository/task/mocks"
	"github.com/rzfhlv/go-task/internal/usecase/task"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/rzfhlv/go-task/pkg/param"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ctxKey string

var (
	idKey ctxKey = "id"
)

func TestTaskCreate(t *testing.T) {
	createRequest := model.Task{
		Title:       "Create Unit Test",
		Description: "for completed task",
		Status:      "todo",
	}

	taskModel := model.Task{
		Title:       "Create Unit Test",
		Description: "for completed task",
		Status:      "todo",
		UserID:      1,
	}

	tests := []struct {
		name       string
		reqContext func(ctx context.Context) context.Context
		mockDeps   func(taskRepository *taskmocks.MockTaskRepository)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(1))
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("Create", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)
					if !ok {
						return false
					}

					return val == int64(1)
				}), mock.MatchedBy(func(task model.Task) bool {
					return task.Title == createRequest.Title &&
						task.Description == createRequest.Description &&
						task.Status == createRequest.Status &&
						task.UserID == int64(1)
				})).Return(taskModel, nil)
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when create task to repository",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(1))
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("Create", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)
					if !ok {
						return false
					}

					return val == int64(1)
				}), mock.MatchedBy(func(task model.Task) bool {
					return task.Title == createRequest.Title &&
						task.Description == createRequest.Description &&
						task.Status == createRequest.Status &&
						task.UserID == int64(1)
				})).Return(model.Task{}, errors.New("some error"))
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get user id from context",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, idKey, int64(1))
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.AssertNotCalled(t, "Create")
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusForbidden, "forbidden access"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepository := taskmocks.MockTaskRepository{}

			ctx := context.Background()
			ctx = tt.reqContext(ctx)

			tt.mockDeps(&taskRepository)

			usecase := task.New(&taskRepository)
			result, err := usecase.Create(ctx, createRequest)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTaskGetByUserID(t *testing.T) {
	userId := int64(1)
	paramReq := param.Param{
		Page:  1,
		Limit: 1,
	}

	tasks := []model.Task{
		{
			ID:          1,
			Title:       "Unit Test",
			Description: "for completness",
			Status:      "todo",
			UserID:      1,
		},
		{
			ID:          2,
			Title:       "Code Review",
			Description: "for completness",
			Status:      "todo",
			UserID:      1,
		},
	}

	tests := []struct {
		name       string
		mockDeps   func(taskRepository *taskmocks.MockTaskRepository)
		wantResult []model.Task
		wantErr    error
	}{
		{
			name: "success",
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByUserID", mock.Anything, mock.MatchedBy(func(requestUserId int64) bool {
					return requestUserId == userId
				}), mock.MatchedBy(func(param param.Param) bool {
					return param.Page == paramReq.Page && param.Limit == paramReq.Limit
				})).Return(tasks, nil)

				taskRepository.On("Count", mock.Anything).Return(int64(2), nil)
			},
			wantResult: tasks,
			wantErr:    nil,
		},
		{
			name: "error when get count task to repository",
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByUserID", mock.Anything, mock.MatchedBy(func(requestUserId int64) bool {
					return requestUserId == userId
				}), mock.MatchedBy(func(param param.Param) bool {
					return param.Page == paramReq.Page && param.Limit == paramReq.Limit
				})).Return(tasks, nil)

				taskRepository.On("Count", mock.Anything).Return(int64(0), errors.New("some error"))
			},
			wantResult: []model.Task{},
			wantErr:    errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "success when task result empty array",
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByUserID", mock.Anything, mock.MatchedBy(func(requestUserId int64) bool {
					return requestUserId == userId
				}), mock.MatchedBy(func(param param.Param) bool {
					return param.Page == paramReq.Page && param.Limit == paramReq.Limit
				})).Return([]model.Task{}, nil)

				taskRepository.On("Count", mock.Anything).Return(int64(2), nil)
			},
			wantResult: []model.Task{},
			wantErr:    nil,
		},
		{
			name: "error when get task by user id",
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByUserID", mock.Anything, mock.MatchedBy(func(requestUserId int64) bool {
					return requestUserId == userId
				}), mock.MatchedBy(func(param param.Param) bool {
					return param.Page == paramReq.Page && param.Limit == paramReq.Limit
				})).Return([]model.Task{}, errors.New("some error"))

				taskRepository.AssertNotCalled(t, "Count")
			},
			wantResult: []model.Task{},
			wantErr:    errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepository := taskmocks.MockTaskRepository{}

			tt.mockDeps(&taskRepository)

			usecase := task.New(&taskRepository)
			result, err := usecase.GetByUserID(context.Background(), userId, &paramReq)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTaskGetByID(t *testing.T) {
	taskId := int64(1)
	userId := int64(1)

	taskModel := model.Task{
		ID:          taskId,
		Title:       "Unit Test",
		Description: "for completness",
		Status:      "todo",
		UserID:      userId,
	}

	tests := []struct {
		name       string
		reqContext func(ctx context.Context) context.Context
		mockDeps   func(taskRepository *taskmocks.MockTaskRepository)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(taskModel, nil)
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when get by id",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(model.Task{}, errors.New("some error"))
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get by id sql no rows",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(model.Task{}, sql.ErrNoRows)
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusNotFound, "task not found"),
		},
		{
			name: "error when get user id from context",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, idKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.AssertNotCalled(t, "GetByID")
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusForbidden, "forbidden access"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepository := taskmocks.MockTaskRepository{}

			ctx := context.Background()
			ctx = tt.reqContext(ctx)

			tt.mockDeps(&taskRepository)

			usecase := task.New(&taskRepository)
			result, err := usecase.GetByID(ctx, taskId)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTaskUpdate(t *testing.T) {
	taskId := int64(1)
	userId := int64(1)

	updateRequest := model.Task{
		ID:          taskId,
		Title:       "Create Unit Test",
		Description: "for completed task",
		Status:      "todo",
	}

	taskModel := model.Task{
		ID:          taskId,
		Title:       "Unit Test",
		Description: "for completness",
		Status:      "todo",
		UserID:      userId,
	}

	tests := []struct {
		name       string
		reqContext func(ctx context.Context) context.Context
		mockDeps   func(taskRepository *taskmocks.MockTaskRepository)
		wantResult model.Task
		wantErr    error
	}{
		{
			name: "success",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(taskModel, nil)

				taskRepository.On("Update", mock.Anything, mock.MatchedBy(func(ts model.Task) bool {
					return ts.Title == updateRequest.Title &&
						ts.Description == updateRequest.Description &&
						ts.Status == updateRequest.Status
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(taskModel, nil)
			},
			wantResult: taskModel,
			wantErr:    nil,
		},
		{
			name: "error when update task",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(taskModel, nil)

				taskRepository.On("Update", mock.Anything, mock.MatchedBy(func(ts model.Task) bool {
					return ts.Title == updateRequest.Title &&
						ts.Description == updateRequest.Description &&
						ts.Status == updateRequest.Status
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(model.Task{}, errors.New("some error"))
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get by id",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(model.Task{}, errors.New("some error"))

				taskRepository.AssertNotCalled(t, "Update")
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get user id from context",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, idKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.AssertNotCalled(t, "GetByID")
				taskRepository.AssertNotCalled(t, "Update")
			},
			wantResult: model.Task{},
			wantErr:    errs.NewErrs(http.StatusForbidden, "forbidden access"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepository := taskmocks.MockTaskRepository{}

			ctx := context.Background()
			ctx = tt.reqContext(ctx)

			tt.mockDeps(&taskRepository)

			usecase := task.New(&taskRepository)
			result, err := usecase.Update(ctx, updateRequest)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTaskDelete(t *testing.T) {
	taskId := int64(1)
	userId := int64(1)

	taskModel := model.Task{
		ID:          taskId,
		Title:       "Unit Test",
		Description: "for completness",
		Status:      "todo",
		UserID:      userId,
	}

	tests := []struct {
		name       string
		reqContext func(ctx context.Context) context.Context
		mockDeps   func(taskRepository *taskmocks.MockTaskRepository)
		wantErr    error
	}{
		{
			name: "success",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(taskModel, nil)

				taskRepository.On("Delete", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "error when delete task",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(taskModel, nil)

				taskRepository.On("Delete", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(errors.New("some error"))
			},
			wantErr: errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get by id",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.On("GetByID", mock.Anything, mock.MatchedBy(func(id int64) bool {
					return id == taskId
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == userId
				})).Return(model.Task{}, errors.New("some error"))

				taskRepository.AssertNotCalled(t, "Delete")
			},
			wantErr: errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get user id from context",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, idKey, userId)
				return ctx
			},
			mockDeps: func(taskRepository *taskmocks.MockTaskRepository) {
				taskRepository.AssertNotCalled(t, "GetByID")
				taskRepository.AssertNotCalled(t, "Delete")
			},
			wantErr: errs.NewErrs(http.StatusForbidden, "forbidden access"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepository := taskmocks.MockTaskRepository{}

			ctx := context.Background()
			ctx = tt.reqContext(ctx)

			tt.mockDeps(&taskRepository)

			usecase := task.New(&taskRepository)
			err := usecase.Delete(ctx, taskId)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
