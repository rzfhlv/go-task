package task

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/repository/task"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/rzfhlv/go-task/pkg/param"
)

type TaskUsecase interface {
	Create(ctx context.Context, task model.Task) (model.Task, error)
	GetByUserID(ctx context.Context, userId int64, param *param.Param) ([]model.Task, error)
	GetByID(ctx context.Context, id int64) (model.Task, error)
	Update(ctx context.Context, task model.Task) (model.Task, error)
	Delete(ctx context.Context, id int64) error
}

type Task struct {
	taskRepository task.TaskRepository
}

func New(taskRepository task.TaskRepository) TaskUsecase {
	return &Task{
		taskRepository: taskRepository,
	}
}

func (t *Task) Create(ctx context.Context, task model.Task) (model.Task, error) {
	val := ctx.Value(auth.IdKey)
	userID, ok := val.(int64)
	if !ok {
		slog.ErrorContext(ctx, "[Usecase.Task] error when get id from context", slog.Any("val", val))
		return model.Task{}, errs.NewErrs(http.StatusForbidden, "forbidden access")
	}

	task.UserID = userID
	result, err := t.taskRepository.Create(ctx, task)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Task] error when call taskRepository.Create", slog.String("error", err.Error()))
		return model.Task{}, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	return result, nil
}

func (t *Task) GetByUserID(ctx context.Context, userId int64, param *param.Param) ([]model.Task, error) {
	result, err := t.taskRepository.GetByUserID(ctx, userId, *param)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Task] error when call taskRepository.GetByUserID", slog.String("error", err.Error()))
		return []model.Task{}, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	if len(result) < 1 {
		result = []model.Task{}
	}

	total, err := t.taskRepository.Count(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Task] error when call taskRepository.Count", slog.String("error", err.Error()))
		return []model.Task{}, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	param.Total = total
	return result, nil
}

func (t *Task) GetByID(ctx context.Context, id int64) (model.Task, error) {
	userId, ok := ctx.Value(auth.IdKey).(int64)
	if !ok {
		slog.ErrorContext(ctx, "[Usecase.Task] error when get user id from context")
		return model.Task{}, errs.NewErrs(http.StatusForbidden, "forbidden access")
	}

	result, err := t.taskRepository.GetByID(ctx, id, userId)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Task] error when call taskRepository.GetByID", slog.String("error", err.Error()))
		if err == sql.ErrNoRows {
			return model.Task{}, errs.NewErrs(http.StatusNotFound, "task not found")
		}

		return model.Task{}, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	return result, nil
}

func (t *Task) Update(ctx context.Context, task model.Task) (model.Task, error) {
	userId, ok := ctx.Value(auth.IdKey).(int64)
	if !ok {
		slog.ErrorContext(ctx, "[Usecase.Task] error when get user id from context")
		return model.Task{}, errs.NewErrs(http.StatusForbidden, "forbidden access")
	}

	check, err := t.GetByID(ctx, task.ID)
	if err != nil {
		return model.Task{}, err
	}

	task.UpdatedAt = time.Now()
	result, err := t.taskRepository.Update(ctx, task, userId)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Task] error when call taskRepository.Update", slog.String("error", err.Error()))
		return model.Task{}, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	result.CreatedAt = check.CreatedAt

	return result, nil
}

func (t *Task) Delete(ctx context.Context, id int64) error {
	userId, ok := ctx.Value(auth.IdKey).(int64)
	if !ok {
		slog.ErrorContext(ctx, "[Usecase.Task] error when get user id from context")
		return errs.NewErrs(http.StatusForbidden, "forbidden access")
	}

	_, err := t.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = t.taskRepository.Delete(ctx, id, userId)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Task] error when call taskRepository.Delete", slog.String("error", err.Error()))
		return errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	return nil
}
