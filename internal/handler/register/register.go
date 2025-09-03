package register

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/usecase/register"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/response/auth"
	"github.com/rzfhlv/go-task/pkg/response/general"
)

type RegisterHandler interface {
	Register(e echo.Context) (err error)
}

type Handler struct {
	usecase register.RegisterUsecase
}

func New(usecase register.RegisterUsecase) RegisterHandler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Register(e echo.Context) (err error) {
	ctx := e.Request().Context()
	user := model.User{}
	err = e.Bind(&user)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, nil)
	}

	err = e.Validate(user)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Register] error when validate the request", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, err.Error()))
	}

	user, jwt, err := h.usecase.Register(ctx, user)
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	resp := auth.AuthResponse{
		JWT:  jwt,
		User: user,
	}

	return e.JSON(http.StatusOK, resp)
}
