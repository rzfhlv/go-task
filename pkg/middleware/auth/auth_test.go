package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
	jwtpkg "github.com/rzfhlv/go-task/pkg/jwt"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cachemocks "github.com/rzfhlv/go-task/internal/repository/cache/mocks"
	jwtmocks "github.com/rzfhlv/go-task/pkg/jwt/mocks"
)

func TestAuthBearer(t *testing.T) {
	tests := []struct {
		name       string
		token      func() string
		Header     string
		TokenType  string
		statusCode int
		mockDeps   func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string)
		wantResult string
	}{
		{
			name: "success",
			token: func() string {
				cfg := config.All("../../../")
				jwtPkg := jwtpkg.New(cfg)
				result, _ := jwtPkg.Generate(model.User{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
				}, "jti-id-1")

				return result.AccessToken
			},
			Header:     "Authorization",
			TokenType:  "Bearer",
			statusCode: http.StatusOK,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.On("ValidateToken", mock.MatchedBy(func(t string) bool {
					return t == token
				})).Return(&jwtpkg.JWTClaim{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
					RegisteredClaims: jwt.RegisteredClaims{
						ID: "jti-id-1",
					},
				}, nil)

				cacheRepository.On("Get", context.Background(), mock.MatchedBy(func(jti string) bool {
					return jti == "jti-id-1"
				})).Return("1", nil)
			},
			wantResult: "test",
		},
		{
			name: "error when claim id not match with cache data",
			token: func() string {
				cfg := config.All("../../../")
				jwtPkg := jwtpkg.New(cfg)
				result, _ := jwtPkg.Generate(model.User{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
				}, "jti-id-1")

				return result.AccessToken
			},
			Header:     "Authorization",
			TokenType:  "Bearer",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.On("ValidateToken", mock.MatchedBy(func(t string) bool {
					return t == token
				})).Return(&jwtpkg.JWTClaim{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
					RegisteredClaims: jwt.RegisteredClaims{
						ID: "jti-id-1",
					},
				}, nil)

				cacheRepository.On("Get", context.Background(), mock.MatchedBy(func(jti string) bool {
					return jti == "jti-id-1"
				})).Return("2", nil)
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
		{
			name: "error when parse cache data not alpha numeric",
			token: func() string {
				cfg := config.All("../../../")
				jwtPkg := jwtpkg.New(cfg)
				result, _ := jwtPkg.Generate(model.User{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
				}, "jti-id-1")

				return result.AccessToken
			},
			Header:     "Authorization",
			TokenType:  "Bearer",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.On("ValidateToken", mock.MatchedBy(func(t string) bool {
					return t == token
				})).Return(&jwtpkg.JWTClaim{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
					RegisteredClaims: jwt.RegisteredClaims{
						ID: "jti-id-1",
					},
				}, nil)

				cacheRepository.On("Get", context.Background(), mock.MatchedBy(func(jti string) bool {
					return jti == "jti-id-1"
				})).Return("ab", nil)
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
		{
			name: "error when get data from cache",
			token: func() string {
				cfg := config.All("../../../")
				jwtPkg := jwtpkg.New(cfg)
				result, _ := jwtPkg.Generate(model.User{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
				}, "jti-id-1")

				return result.AccessToken
			},
			Header:     "Authorization",
			TokenType:  "Bearer",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.On("ValidateToken", mock.MatchedBy(func(t string) bool {
					return t == token
				})).Return(&jwtpkg.JWTClaim{
					ID:    1,
					Name:  "John",
					Email: "john@mail.com",
					RegisteredClaims: jwt.RegisteredClaims{
						ID: "jti-id-1",
					},
				}, nil)

				cacheRepository.On("Get", context.Background(), mock.MatchedBy(func(jti string) bool {
					return jti == "jti-id-1"
				})).Return("", errors.New("some error"))
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
		{
			name: "error when validate token",
			token: func() string {
				return "invalidtoken"
			},
			Header:     "Authorization",
			TokenType:  "Bearer",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.On("ValidateToken", mock.MatchedBy(func(t string) bool {
					return t == token
				})).Return(nil, errors.New("some error"))

				cacheRepository.AssertNotCalled(t, "Get")
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
		{
			name: "error when token is empty",
			token: func() string {
				return ""
			},
			Header:     "Authorization",
			TokenType:  "Bearer",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.AssertNotCalled(t, "ValidateToken")
				cacheRepository.AssertNotCalled(t, "Get")
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
		{
			name: "error when token type not bearer",
			token: func() string {
				return ""
			},
			Header:     "Authorization",
			TokenType:  "Beer ",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.AssertNotCalled(t, "ValidateToken")
				cacheRepository.AssertNotCalled(t, "Get")
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
		{
			name: "error when value header less than 2",
			token: func() string {
				return ""
			},
			Header:     "x",
			TokenType:  "",
			statusCode: http.StatusUnauthorized,
			mockDeps: func(cacheRepository *cachemocks.MockCacheRepository, jwtImpl *jwtmocks.MockJWTInterface, token string) {
				jwtImpl.AssertNotCalled(t, "ValidateToken")
				cacheRepository.AssertNotCalled(t, "Get")
			},
			wantResult: "{\"success\":false,\"error\":\"unauthorized\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token()
			cacheRepository := cachemocks.MockCacheRepository{}
			jwtImpl := jwtmocks.MockJWTInterface{}

			tt.mockDeps(&cacheRepository, &jwtImpl, token)

			e := echo.New()
			auth := auth.New(&cacheRepository, &jwtImpl)
			e.Use(auth.Bearer)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(tt.Header, tt.TokenType+" "+token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := func(c echo.Context) error {
				return c.String(http.StatusOK, tt.wantResult)
			}
			auth.Bearer(handler)(c)

			// assert.NoError(t, err)
			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, tt.wantResult, rec.Body.String())
		})
	}
}
