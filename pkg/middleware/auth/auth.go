package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rzfhlv/go-task/internal/repository/cache"
	"github.com/rzfhlv/go-task/pkg/jwt"
	"github.com/rzfhlv/go-task/pkg/response/general"
)

type ctxKey string

var (
	IdKey  ctxKey = "id"
	JtiKey ctxKey = "jti"
)

type AuthMiddleware interface {
	Bearer(next echo.HandlerFunc) echo.HandlerFunc
}

type Auth struct {
	cacheRepository cache.CacheRepository
	jwt             jwt.JWTInterface
}

func New(cacheRepository cache.CacheRepository, jwt jwt.JWTInterface) AuthMiddleware {
	return &Auth{
		cacheRepository: cacheRepository,
		jwt:             jwt,
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

		jti := claims.RegisteredClaims.ID
		val, err := a.cacheRepository.Get(context.Background(), jti)
		if err != nil {
			slog.Error("error when get token from cahce", slog.String("error", err.Error()))
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		valInt, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			slog.Error("error when parse value from cache", slog.String("error", err.Error()))
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		if claims.ID != valInt {
			slog.Error("error id not match", slog.Any("claims_id", claims.ID), slog.Any("val_from_cache", valInt))
			return c.JSON(http.StatusUnauthorized, general.Set(false, nil, nil, nil, "unauthorized"))
		}

		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, IdKey, valInt)
		ctx = context.WithValue(ctx, JtiKey, jti)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
