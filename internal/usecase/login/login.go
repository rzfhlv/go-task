package login

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/repository/cache"
	"github.com/rzfhlv/go-task/internal/repository/user"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/hasher"
	"github.com/rzfhlv/go-task/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase interface {
	Login(ctx context.Context, login model.Login) (model.User, model.JWT, error)
}

type Login struct {
	userRepository  user.UserRepository
	cacheRepository cache.CacheRepository
	hasher          hasher.HashPassword
	jwt             jwt.JWTInterface
}

func New(userRepository user.UserRepository, cacheRepository cache.CacheRepository, hasher hasher.HashPassword, jwt jwt.JWTInterface) LoginUsecase {
	return &Login{
		userRepository:  userRepository,
		cacheRepository: cacheRepository,
		hasher:          hasher,
		jwt:             jwt,
	}
}

func (l *Login) Login(ctx context.Context, login model.Login) (model.User, model.JWT, error) {
	result := model.User{}
	jwt := model.JWT{}

	user, err := l.userRepository.GetByEmail(ctx, login.Email)
	if err != nil {
		slog.InfoContext(ctx, "[Usecase.Login] error when call userRepository.GetByEmail", slog.String("err", err.Error()))
		if err == sql.ErrNoRows {
			return result, jwt, errs.NewErrs(http.StatusUnauthorized, "unauthorized")
		}

		return result, jwt, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	err = l.hasher.VerifyPassword(user.Password, login.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		slog.InfoContext(ctx, "[Usecase.Login] error when call hasher.VerifyPassword", slog.String("err", err.Error()))
		return result, jwt, errs.NewErrs(http.StatusUnauthorized, "invalid credentials")
	}

	jti := uuid.NewString()
	token, err := l.jwt.Generate(user, jti)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Login] error when call jwt.Generate", slog.String("error", err.Error()))
		return result, jwt, errs.NewErrs(http.StatusUnauthorized, "invalid credentials")
	}

	cfg := config.Get()
	err = l.cacheRepository.Set(ctx, jti, user.ID, cfg.JWT.ExpiresIn)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Login] error when call redis.Set", slog.String("error", err.Error()))
		return result, jwt, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	return user, token, nil
}
