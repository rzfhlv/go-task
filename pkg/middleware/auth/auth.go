package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rzfhlv/go-task/pkg/jwt"
	"github.com/rzfhlv/go-task/pkg/response/general"
)

type ctxKey string

var (
	IdKey ctxKey = "id"
)

type AuthMiddleware interface {
	Bearer(next echo.HandlerFunc) echo.HandlerFunc
}

type Auth struct {
	redis *redis.Client
	jwt   jwt.JWTInterface
}

func New(redis *redis.Client, jwt jwt.JWTInterface) AuthMiddleware {
	return &Auth{
		redis: redis,
		jwt:   jwt,
	}
}

func (a *Auth) Bearer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		split := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(split) < 2 {
			slog.Error("missing header authorization")
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		if split[0] != "Bearer" {
			slog.Error("missing bearer authorization")
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		if split[1] == "" {
			slog.Error("missing bearer token")
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		claims, err := a.jwt.ValidateToken(split[1])
		if err != nil {
			slog.Error("error when validate token", slog.String("error", err.Error()))
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		err = a.redis.Get(context.Background(), claims.RegisteredClaims.ID).Err()
		if err != nil {
			slog.Error("error when get token from cahce", slog.String("error", err.Error()))
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		slog.Info("id from validate", slog.Any("id", claims.ID))
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, IdKey, int64(claims.ID))
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
