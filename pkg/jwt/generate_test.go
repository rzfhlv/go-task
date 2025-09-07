package jwt_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
	jwtpkg "github.com/rzfhlv/go-task/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

var (
	jwtPkg jwtpkg.JWTInterface
)

func TestJwtGenerate(t *testing.T) {
	userModel := model.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@mail.com",
	}

	tests := []struct {
		name      string
		user      model.User
		jti       string
		secret    string
		expiresIn time.Duration
		wantErr   error
	}{
		{
			name:      "success",
			user:      userModel,
			jti:       "jti-id-1",
			secret:    "verysecret",
			expiresIn: time.Duration(5 * time.Minute),
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Configuration{
				JWT: config.JWTConfiguration{
					Secret:    tt.secret,
					ExpiresIn: tt.expiresIn,
				},
			}

			jwtPkg = jwtpkg.New(cfg)
			result, err := jwtPkg.Generate(tt.user, tt.jti)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantErr, err)

			if tt.wantErr == nil {
				token, err := jwt.ParseWithClaims(result.AccessToken, &jwtpkg.JWTClaim{}, func(t *jwt.Token) (any, error) {
					return []byte(cfg.JWT.Secret), nil
				})

				assert.NoError(t, err)

				claims, ok := token.Claims.(*jwtpkg.JWTClaim)
				assert.True(t, ok)
				assert.Equal(t, tt.user.ID, claims.ID)
				assert.Equal(t, tt.user.Email, claims.Email)
				assert.Equal(t, tt.user.Name, claims.Name)
				assert.Equal(t, tt.jti, claims.RegisteredClaims.ID)
			}
		})
	}
}
