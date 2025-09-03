package register

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/repository/user"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/hasher"
	"github.com/rzfhlv/go-task/pkg/jwt"
)

type RegisterUsecase interface {
	Register(ctx context.Context, user model.User) (model.User, model.JWT, error)
}

type Register struct {
	userRepository user.UserRepository
	hasher         hasher.HashPassword
	jwt            jwt.JWTInterface
}

func New(userRepository user.UserRepository, hasher hasher.HashPassword, jwt jwt.JWTInterface) RegisterUsecase {
	return &Register{
		userRepository: userRepository,
		hasher:         hasher,
		jwt:            jwt,
	}
}

func (r *Register) Register(ctx context.Context, user model.User) (model.User, model.JWT, error) {
	result := model.User{}
	jwt := model.JWT{}

	hashPassword, err := r.hasher.HashedPassword(user.Password)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Register] error when call hasher.HashedPassword", slog.String("error", err.Error()))
		return result, jwt, errs.NewErrs(http.StatusUnprocessableEntity, "hasher error")
	}

	user.Password = hashPassword

	_, err = r.userRepository.GetByEmail(ctx, user.Email)
	if err == nil {
		slog.InfoContext(ctx, "[Usecase.Register] duplicate data when call userRepository.GetByEmail", slog.String("email", user.Email))
		return result, jwt, errs.NewErrs(http.StatusUnprocessableEntity, "email already exists")
	}

	user.CreatedAt = time.Now()
	result, err = r.userRepository.Create(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Register] error when call userRepository.Create", slog.String("error", err.Error()))
		return result, jwt, errs.NewErrs(http.StatusInternalServerError, "something went wrong")
	}

	jti := uuid.NewString()
	token, err := r.jwt.Generate(result, jti)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Register] error when call jwt.Generate", slog.String("error", err.Error()))
		return result, jwt, errs.NewErrs(http.StatusUnprocessableEntity, "failed generated token")
	}

	return result, token, nil
}
