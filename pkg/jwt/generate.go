package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rzfhlv/go-task/internal/model"
)

const bearer = "Bearer"

type JWTClaim struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (j *JWTImpl) Generate(user model.User, jti string) (jwtModel model.JWT, err error) {
	expirationTime := time.Now().Add(j.cfg.JWT.ExpiresIn)
	claims := &JWTClaim{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.App.Name,
			Subject:   user.Name,
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.cfg.JWT.Secret))
	if err != nil {
		return
	}

	jwtModel.AccessToken = tokenString
	jwtModel.TokenType = bearer
	jwtModel.ExpiresIn = int(j.cfg.JWT.ExpiresIn.Seconds())

	return
}
