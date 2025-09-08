package task_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/handler/task"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/presenter/rest"
	taskmocks "github.com/rzfhlv/go-task/internal/usecase/task/mocks"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/rzfhlv/go-task/pkg/param"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	createRequest = model.Task{
		Title:       "Task 1",
		Description: "for test",
		Status:      "todo",
	}

	taskModel = model.Task{
		ID:          1,
		Title:       "Task 1",
		Description: "for test",
		Status:      "todo",
		UserID:      1,
	}
)

func TestHandlerTaskCreate(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    string
		mockDeps   func(taskUsecase *taskmocks.MockTaskUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name:    "success",
			reqBody: `{"title": "Task 1", "description": "for test", "status": "todo"}`,
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Create", context.Background(), createRequest).
					Return(taskModel, nil)
			},
			statusCode: http.StatusCreated,
			wantErr:    nil,
		},
		{
			name:    "error when call login usecase",
			reqBody: `{"title": "Task 1", "description": "for test", "status": "todo"}`,
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Create", context.Background(), createRequest).
					Return(model.Task{}, errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name:    "error when call login usecase with custome error message",
			reqBody: `{"title": "Task 1", "description": "for test", "status": "todo"}`,
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Create", context.Background(), createRequest).
					Return(model.Task{}, errs.NewErrs(http.StatusForbidden, "forbidden access"))
			},
			statusCode: http.StatusForbidden,
			wantErr:    nil,
		},
		{
			name:    "error when validate request",
			reqBody: `{"title": "", "description": "for test", "status": "todo"}`,
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Create", context.Background(), createRequest).
					Return(model.Task{}, errors.New("title is required"))
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
		{
			name:    "error when binding request",
			reqBody: `{`,
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Create", context.Background(), createRequest).
					Return(model.Task{}, errors.New("error binding request"))
			},
			statusCode: http.StatusUnprocessableEntity,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskUsecase := taskmocks.MockTaskUsecase{}

			tt.mockDeps(&taskUsecase)

			handler := task.New(&taskUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/v1/task", strings.NewReader(tt.reqBody))
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := handler.Create(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHandlerTaskGetByUserID(t *testing.T) {
	tests := []struct {
		name       string
		reqParam   string
		mockCtx    func(ctx context.Context) context.Context
		mockDeps   func(taskUsecase *taskmocks.MockTaskUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name:     "success",
			reqParam: "?page=2",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("GetByUserID", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == taskModel.UserID
				}), mock.MatchedBy(func(p *param.Param) bool {
					return p.Page == 2
				})).
					Return([]model.Task{taskModel}, nil)
			},
			statusCode: http.StatusOK,
			wantErr:    nil,
		},
		{
			name:     "error when call login usecase",
			reqParam: "?page=2",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("GetByUserID", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == taskModel.UserID
				}), mock.MatchedBy(func(p *param.Param) bool {
					return p.Page == 2
				})).
					Return([]model.Task{taskModel}, errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name:     "error when call login usecase with custome error message",
			reqParam: "?page=2",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("GetByUserID", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(uid int64) bool {
					return uid == taskModel.UserID
				}), mock.MatchedBy(func(p *param.Param) bool {
					return p.Page == 2
				})).
					Return([]model.Task{taskModel}, errs.NewErrs(http.StatusForbidden, "forbidden access"))
			},
			statusCode: http.StatusForbidden,
			wantErr:    nil,
		},
		{
			name:     "error when binding request param",
			reqParam: "?page=satu",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "GetByUserID")
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
		{
			name:     "error when get user id from context",
			reqParam: "?page=satu",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "GetByUserID")
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskUsecase := taskmocks.MockTaskUsecase{}

			tt.mockDeps(&taskUsecase)

			handler := task.New(&taskUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodGet, "/v1/task"+tt.reqParam, nil)
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			baseContext := ctx.Request().Context()
			baseContext = tt.mockCtx(baseContext)
			ctx.SetRequest(ctx.Request().WithContext(baseContext))

			err := handler.GetByUserID(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHandlerTaskGetByID(t *testing.T) {
	tests := []struct {
		name       string
		pathParam  string
		mockCtx    func(ctx context.Context) context.Context
		mockDeps   func(taskUsecase *taskmocks.MockTaskUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name:      "success",
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(id int64) bool {
					return id == taskModel.ID
				})).
					Return(taskModel, nil)
			},
			statusCode: http.StatusOK,
			wantErr:    nil,
		},
		{
			name:      "error when call login usecase",
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(id int64) bool {
					return id == taskModel.ID
				})).
					Return(model.Task{}, errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name:      "error when call login usecase with custome error message",
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(id int64) bool {
					return id == taskModel.ID
				})).
					Return(model.Task{}, errs.NewErrs(http.StatusForbidden, "forbidden access"))
			},
			statusCode: http.StatusForbidden,
			wantErr:    nil,
		},
		{
			name:      "error when parse request path param",
			pathParam: "satu",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "GetByUserID")
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskUsecase := taskmocks.MockTaskUsecase{}

			tt.mockDeps(&taskUsecase)

			handler := task.New(&taskUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodGet, "/v1/task/"+tt.pathParam, nil)
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			baseContext := ctx.Request().Context()
			baseContext = tt.mockCtx(baseContext)
			ctx.SetRequest(ctx.Request().WithContext(baseContext))
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.pathParam)

			err := handler.GetByID(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHandlerTaskUpdate(t *testing.T) {
	updatedRequest := model.Task{
		ID:          1,
		Title:       "Task 1",
		Description: "for test",
		Status:      "completed",
	}

	taskModelCompleted := model.Task{
		ID:          1,
		Title:       "Task 1",
		Description: "for test",
		Status:      "completed",
		UserID:      1,
	}

	tests := []struct {
		name       string
		reqBody    string
		pathParam  string
		mockCtx    func(ctx context.Context) context.Context
		mockDeps   func(taskUsecase *taskmocks.MockTaskUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name:      "success",
			reqBody:   `{"title": "Task 1", "description": "for test", "status": "completed"}`,
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), updatedRequest).
					Return(taskModelCompleted, nil)
			},
			statusCode: http.StatusOK,
			wantErr:    nil,
		},
		{
			name:      "error when call update usecase",
			reqBody:   `{"title": "Task 1", "description": "for test", "status": "completed"}`,
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), updatedRequest).
					Return(model.Task{}, errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name:      "error when call update usecase with custom error",
			reqBody:   `{"title": "Task 1", "description": "for test", "status": "completed"}`,
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Update", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), updatedRequest).
					Return(model.Task{}, errs.NewErrs(http.StatusForbidden, "forbidden access"))
			},
			statusCode: http.StatusForbidden,
			wantErr:    nil,
		},
		{
			name:      "error when validate request",
			reqBody:   `{"title": "", "description": "for test", "status": "completed"}`,
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "Update")
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
		{
			name:      "error when binding request",
			reqBody:   `{`,
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "Update")
			},
			statusCode: http.StatusUnprocessableEntity,
			wantErr:    nil,
		},
		{
			name:      "error when binding request",
			reqBody:   `{"title": "Task 1", "description": "for test", "status": "completed"}`,
			pathParam: "satu",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "Update")
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskUsecase := taskmocks.MockTaskUsecase{}

			tt.mockDeps(&taskUsecase)

			handler := task.New(&taskUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPut, "/v1/task/"+tt.pathParam, strings.NewReader(tt.reqBody))
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			baseContext := ctx.Request().Context()
			baseContext = tt.mockCtx(baseContext)
			ctx.SetRequest(ctx.Request().WithContext(baseContext))
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.pathParam)

			err := handler.Update(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHandlerTaskDelete(t *testing.T) {
	tests := []struct {
		name       string
		pathParam  string
		mockCtx    func(ctx context.Context) context.Context
		mockDeps   func(taskUsecase *taskmocks.MockTaskUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name:      "success",
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Delete", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(id int64) bool {
					return id == taskModel.ID
				})).
					Return(nil)
			},
			statusCode: http.StatusOK,
			wantErr:    nil,
		},
		{
			name:      "error when call login usecase",
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Delete", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(id int64) bool {
					return id == taskModel.ID
				})).
					Return(errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name:      "error when call login usecase with custome error message",
			pathParam: "1",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.On("Delete", mock.MatchedBy(func(ctx context.Context) bool {
					val, ok := ctx.Value(auth.IdKey).(int64)

					return val == taskModel.UserID && ok
				}), mock.MatchedBy(func(id int64) bool {
					return id == taskModel.ID
				})).
					Return(errs.NewErrs(http.StatusForbidden, "forbidden access"))
			},
			statusCode: http.StatusForbidden,
			wantErr:    nil,
		},
		{
			name:      "error when parse request path param",
			pathParam: "satu",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.IdKey, int64(taskModel.UserID))
				return ctx
			},
			mockDeps: func(taskUsecase *taskmocks.MockTaskUsecase) {
				taskUsecase.AssertNotCalled(t, "GetByUserID")
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskUsecase := taskmocks.MockTaskUsecase{}

			tt.mockDeps(&taskUsecase)

			handler := task.New(&taskUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodDelete, "/v1/task/"+tt.pathParam, nil)
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			baseContext := ctx.Request().Context()
			baseContext = tt.mockCtx(baseContext)
			ctx.SetRequest(ctx.Request().WithContext(baseContext))
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.pathParam)

			err := handler.Delete(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
