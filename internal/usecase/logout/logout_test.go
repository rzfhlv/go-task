package logout_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	cachemocks "github.com/rzfhlv/go-task/internal/repository/cache/mocks"
	"github.com/rzfhlv/go-task/internal/usecase/logout"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ctxKey string

var (
	customKey ctxKey = "jti"
)

func TestLogout(t *testing.T) {
	tests := []struct {
		name       string
		reqContext func(ctx context.Context) context.Context
		mockDeps   func(cacheRepository *cachemocks.MockCacheRepository)
		wantErr    error
	}{
		{
			name: "success",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, "jti-id")
				return ctx
			},
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository) {
				cacheRepository.On("Del", mock.Anything, mock.Anything).Return(int64(1), nil)
			},
			wantErr: nil,
		},
		{
			name: "error when no deleted data",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, "jti-id")
				return ctx
			},
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository) {
				cacheRepository.On("Del", mock.Anything, mock.Anything).Return(int64(0), nil)
			},
			wantErr: errs.NewErrs(http.StatusForbidden, "forbidden access"),
		},
		{
			name: "error when deleted the data",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, auth.JtiKey, "jti-id")
				return ctx
			},
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository) {
				cacheRepository.On("Del", mock.Anything, mock.Anything).Return(int64(0), errors.New("some error"))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "error when get value from context",
			reqContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, customKey, "jti-id")
				return ctx
			},
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository) {
				cacheRepository.AssertNotCalled(t, "Del")
			},
			wantErr: errs.NewErrs(http.StatusForbidden, "forbidden access"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheRepository := cachemocks.MockCacheRepository{}

			ctx := context.Background()
			ctx = tt.reqContext(ctx)

			tt.mockDeps(&cacheRepository)

			usecase := logout.New(&cacheRepository)
			err := usecase.Logout(ctx)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
