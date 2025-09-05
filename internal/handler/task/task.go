package task

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/usecase/task"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/rzfhlv/go-task/pkg/param"
	"github.com/rzfhlv/go-task/pkg/response/general"
)

type TaskHandler interface {
	Create(e echo.Context) (err error)
	GetByUserID(e echo.Context) (err error)
	GetByID(e echo.Context) (err error)
	Update(e echo.Context) (err error)
	Delete(e echo.Context) (err error)
}

type Handler struct {
	usecase task.TaskUsecase
}

func New(usecase task.TaskUsecase) TaskHandler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Create(e echo.Context) (err error) {
	ctx := e.Request().Context()
	task := model.Task{}
	err = e.Bind(&task)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, general.Set(false, nil, nil, nil, "invalid json"))
	}

	err = e.Validate(task)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Task] error when validate the request", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, err.Error()))
	}

	result, err := h.usecase.Create(ctx, task)
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	msg := "created success"
	return e.JSON(http.StatusCreated, general.Set(true, &msg, nil, result, nil))
}

func (h *Handler) GetByUserID(e echo.Context) (err error) {
	ctx := e.Request().Context()
	param := param.Param{}
	param.Limit = 10
	param.Page = 1

	userId, ok := ctx.Value(auth.IdKey).(int64)
	if !ok {
		slog.ErrorContext(ctx, "[Usecase.Task] error when get id from context")
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, "missing user id in context"))
	}

	err = (&echo.DefaultBinder{}).BindQueryParams(e, &param)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Task] error when bind query param request", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, "invalid query param request"))
	}

	result, err := h.usecase.GetByUserID(ctx, userId, &param)
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	msg := "get data success"
	meta := general.BuildMeta(param, len(result))
	return e.JSON(http.StatusOK, general.Set(true, &msg, meta, result, nil))
}

func (h *Handler) GetByID(e echo.Context) (err error) {
	ctx := e.Request().Context()

	id := e.Param("id")
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Task] error when convert id param to int", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, "invalid path param id"))
	}

	result, err := h.usecase.GetByID(ctx, taskId)
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	msg := "get data success"
	return e.JSON(http.StatusOK, general.Set(true, &msg, nil, result, nil))
}

func (h *Handler) Update(e echo.Context) (err error) {
	ctx := e.Request().Context()

	id := e.Param("id")
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Task] error when convert id param to int", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, "invalid path param id"))
	}

	task := model.Task{}
	err = e.Bind(&task)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, general.Set(false, nil, nil, nil, "invalid json"))
	}

	task.ID = taskId
	result, err := h.usecase.Update(ctx, task)
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	msg := "update data success"
	return e.JSON(http.StatusOK, general.Set(true, &msg, nil, result, nil))
}

func (h *Handler) Delete(e echo.Context) (err error) {
	ctx := e.Request().Context()

	id := e.Param("id")
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Task] error when convert id param to int", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, "invalid path param id"))
	}

	err = h.usecase.Delete(ctx, taskId)
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	msg := "delete data success"
	return e.JSON(http.StatusOK, general.Set(true, &msg, nil, nil, nil))
}
