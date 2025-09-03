package jwt

import (
	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
)

type JWTInterface interface {
	Generate(user model.User, jti string) (jwtModel model.JWT, err error)
	ValidateToken(signedToken string) (claims *JWTClaim, err error)
}

type JWTImpl struct {
	cfg *config.Configuration
}

func New(cfg *config.Configuration) JWTInterface {
	return &JWTImpl{
		cfg: cfg,
	}
}
