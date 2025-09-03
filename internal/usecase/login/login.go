package login

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/rzfhlv/go-task/internal/model"
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
	userRepository user.UserRepository
	hasher         hasher.HashPassword
	jwt            jwt.JWTInterface
}

func New(userRepository user.UserRepository, hasher hasher.HashPassword, jwt jwt.JWTInterface) LoginUsecase {
	return &Login{
		userRepository: userRepository,
		hasher:         hasher,
		jwt:            jwt,
	}
}

func (l *Login) Login(ctx context.Context, login model.Login) (model.User, model.JWT, error) {
	result := model.User{}
	jwt := model.JWT{}

	user, err := l.userRepository.GetByEmail(ctx, login.Email)
	if err != nil {
		slog.InfoContext(ctx, "[Usecase.Login] error when call userRepository.GetByEmail", slog.String("err", err.Error()))
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

	return user, token, nil
}
