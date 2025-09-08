package login_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/handler/login"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/presenter/rest"
	loginmocks "github.com/rzfhlv/go-task/internal/usecase/login/mocks"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestHandlerLogin(t *testing.T) {
	loginRequest := model.Login{
		Email:    "john@mail.com",
		Password: "password",
	}

	userModel := model.User{
		ID:       1,
		Name:     "John",
		Email:    "john@mail.com",
		Password: "password",
	}

	jwtModel := model.JWT{
		AccessToken: "thisistoken",
		TokenType:   "Bearer",
		ExpiresIn:   900,
	}

	tests := []struct {
		name       string
		reqBody    string
		mockDeps   func(loginUsecase *loginmocks.MockLoginUsecase)
		statusCode int
		wantErr    error
	}{
		{
			name:    "success",
			reqBody: `{"email": "john@mail.com", "password": "password"}`,
			mockDeps: func(loginUsecase *loginmocks.MockLoginUsecase) {
				loginUsecase.On("Login", context.Background(), loginRequest).
					Return(userModel, jwtModel, nil)
			},
			statusCode: http.StatusOK,
			wantErr:    nil,
		},
		{
			name:    "error when call login usecase",
			reqBody: `{"email": "john@mail.com", "password": "password"}`,
			mockDeps: func(loginUsecase *loginmocks.MockLoginUsecase) {
				loginUsecase.On("Login", context.Background(), loginRequest).
					Return(model.User{}, model.JWT{}, errors.New("some error"))
			},
			statusCode: http.StatusInternalServerError,
			wantErr:    nil,
		},
		{
			name:    "error when call login usecase with custome error message",
			reqBody: `{"email": "john@mail.com", "password": "password"}`,
			mockDeps: func(loginUsecase *loginmocks.MockLoginUsecase) {
				loginUsecase.On("Login", context.Background(), loginRequest).
					Return(model.User{}, model.JWT{}, errs.NewErrs(http.StatusUnauthorized, "unauthorized"))
			},
			statusCode: http.StatusUnauthorized,
			wantErr:    nil,
		},
		{
			name:    "error when validate request",
			reqBody: `{"email": "john@mail.com", "password": ""}`,
			mockDeps: func(loginUsecase *loginmocks.MockLoginUsecase) {
				loginUsecase.On("Login", context.Background(), loginRequest).
					Return(model.User{}, model.JWT{}, errors.New("password is required"))
			},
			statusCode: http.StatusBadRequest,
			wantErr:    nil,
		},
		{
			name:    "error when binding request",
			reqBody: `{`,
			mockDeps: func(loginUsecase *loginmocks.MockLoginUsecase) {
				loginUsecase.On("Login", context.Background(), loginRequest).
					Return(model.User{}, model.JWT{}, errors.New("error binding request"))
			},
			statusCode: http.StatusUnprocessableEntity,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginUsecase := loginmocks.MockLoginUsecase{}

			tt.mockDeps(&loginUsecase)

			handler := login.New(&loginUsecase)

			e := echo.New()
			e.Validator = &rest.CustomValidator{Validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/v1/login", strings.NewReader(tt.reqBody))
			req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			err := handler.Login(ctx)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
