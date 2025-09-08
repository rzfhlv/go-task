package logout_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/handler/logout"
	"github.com/rzfhlv/go-task/internal/presenter/rest"
	logoutmocks "github.com/rzfhlv/go-task/internal/usecase/logout/mocks"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerLogout(t *testing.T) {
	tests := []struct {
		name       string
		mockCtx    func(ctx context.Context) context.Context
		mockDeps   func(logoutUsecase *logoutmocks.MockLogoutUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name: "success",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, "some-id")
				return ctx
			},
			mockDeps: func(logoutUsecase *logoutmocks.MockLogoutUsecase) {
				logoutUsecase.On("Logout", mock.MatchedBy(func(ctx context.Context) bool {
					key, ok := ctx.Value(auth.JtiKey).(string)

					return key == "some-id" && ok
				})).Return(nil)
			},
			statusCode: http.StatusOK,
			wantErr:    nil,
		},
		{
			name: "error when logout usecase",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, "some-id")
				return ctx
			},
			mockDeps: func(logoutUsecase *logoutmocks.MockLogoutUsecase) {
				logoutUsecase.On("Logout", mock.MatchedBy(func(ctx context.Context) bool {
					key, ok := ctx.Value(auth.JtiKey).(string)

					return key == "some-id" && ok
				})).Return(errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name: "error when logout usecase",
			mockCtx: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, "some-id")
				return ctx
			},
			mockDeps: func(logoutUsecase *logoutmocks.MockLogoutUsecase) {
				logoutUsecase.On("Logout", mock.MatchedBy(func(ctx context.Context) bool {
					key, ok := ctx.Value(auth.JtiKey).(string)

					return key == "some-id" && ok
				})).Return(errs.NewErrs(http.StatusForbidden, "forbidden access"))
			},
			statusCode: http.StatusForbidden,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logoutUsecase := logoutmocks.MockLogoutUsecase{}

			tt.mockDeps(&logoutUsecase)

			handler := logout.New(&logoutUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/v1/logout", nil)
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			baseContext := ctx.Request().Context()
			baseContext = tt.mockCtx(baseContext)
			ctx.SetRequest(ctx.Request().WithContext(baseContext))

			err := handler.Logout(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
