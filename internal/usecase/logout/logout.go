package logout

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rzfhlv/go-task/internal/repository/cache"
	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
)

type LogoutUsecase interface {
	Logout(ctx context.Context) error
}

type Logout struct {
	cacheRepository cache.CacheRepository
}

func New(cacheRepository cache.CacheRepository) LogoutUsecase {
	return &Logout{
		cacheRepository: cacheRepository,
	}
}

func (l *Logout) Logout(ctx context.Context) error {
	key, ok := ctx.Value(auth.JtiKey).(string)
	if !ok {
		slog.ErrorContext(ctx, "[Usecase.Logout] error when get jti id from context")
		return errs.NewErrs(http.StatusForbidden, "forbidden access")
	}

	deleted, err := l.cacheRepository.Del(ctx, key)
	if err != nil {
		slog.ErrorContext(ctx, "[Usecase.Logout] error when call cacheRepository.Del", slog.String("error", err.Error()))
		return err
	}

	if deleted < 1 {
		slog.ErrorContext(ctx, "[Usecase.Logout] no data deleted from cache", slog.String("key", key))
		return errs.NewErrs(http.StatusForbidden, "forbidden access")
	}

	return nil
}
