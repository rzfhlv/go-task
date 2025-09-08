package login

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/usecase/login"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/response/auth"
	"github.com/rzfhlv/go-task/pkg/response/general"
)

type LoginHandler interface {
	Login(e echo.Context) (err error)
}

type Handler struct {
	usecase login.LoginUsecase
}

func New(usecase login.LoginUsecase) LoginHandler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Login(e echo.Context) (err error) {
	ctx := e.Request().Context()
	login := model.Login{}
	err = e.Bind(&login)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, general.Set(false, nil, nil, nil, "error when binding request"))
	}

	err = e.Validate(login)
	if err != nil {
		slog.ErrorContext(ctx, "[Handler.Login] error when validate the request", slog.String("error", err.Error()))
		return e.JSON(http.StatusBadRequest, general.Set(false, nil, nil, nil, err.Error()))
	}

	user, jwt, err := h.usecase.Login(ctx, login)
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
