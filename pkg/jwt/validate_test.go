package jwt_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
	jwtpkg "github.com/rzfhlv/go-task/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJwtValidate(t *testing.T) {
	tests := []struct {
		name    string
		token   func() string
		wantErr error
	}{
		{
			name: "success",
			token: func() string {
				cfg := &config.Configuration{
					JWT: config.JWTConfiguration{
						Secret:    "verysecret",
						ExpiresIn: time.Duration(5 * time.Minute),
					},
				}

				user := model.User{
					ID:    1,
					Name:  "John Doe",
					Email: "john@mail.com",
				}

				jwtRequest := jwtpkg.New(cfg)
				jwtResponse, _ := jwtRequest.Generate(user, "jti-id-1")
				return jwtResponse.AccessToken
			},
			wantErr: nil,
		},
		{
			name: "error when validate jwt token",
			token: func() string {
				return "invalidtoken"
			},
			wantErr: errors.New("token is malformed: token contains an invalid number of segments"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Configuration{
				JWT: config.JWTConfiguration{
					Secret:    "verysecret",
					ExpiresIn: time.Duration(5 * time.Minute),
				},
			}

			token := tt.token()
			jwtPkg := jwtpkg.New(cfg)
			_, err := jwtPkg.ValidateToken(token)

			if err != nil {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
