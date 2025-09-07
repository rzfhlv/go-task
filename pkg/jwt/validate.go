package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

func (j *JWTImpl) ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(j.cfg.JWT.Secret), nil
		},
	)
	if err != nil {
		return
	}

	claims, _ = token.Claims.(*JWTClaim)
	return
}
