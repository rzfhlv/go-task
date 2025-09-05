package logout

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/usecase/logout"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/response/general"
)

type LogoutHandler interface {
	Logout(e echo.Context) (err error)
}

type Handler struct {
	usecase logout.LogoutUsecase
}

func New(usecase logout.LogoutUsecase) LogoutHandler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) Logout(e echo.Context) (err error) {
	err = h.usecase.Logout(e.Request().Context())
	if err != nil {
		if httpErr, ok := err.(*errs.HttpError); ok {
			return e.JSON(httpErr.StatusCode, general.Set(false, nil, nil, nil, httpErr.Message))
		}

		return e.JSON(http.StatusInternalServerError, general.Set(false, nil, nil, nil, "something went wrong"))
	}

	msg := "logout successful"
	return e.JSON(http.StatusOK, general.Set(true, &msg, nil, nil, nil))
}
