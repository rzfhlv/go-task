package rest

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rzfhlv/go-task/config"
	loginhandler "github.com/rzfhlv/go-task/internal/handler/login"
	registerhandler "github.com/rzfhlv/go-task/internal/handler/register"
	"github.com/rzfhlv/go-task/internal/infrastructure"
	"github.com/rzfhlv/go-task/internal/repository/user"
	"github.com/rzfhlv/go-task/internal/usecase/login"
	"github.com/rzfhlv/go-task/internal/usecase/register"
	"github.com/rzfhlv/go-task/pkg/hasher"
	"github.com/rzfhlv/go-task/pkg/jwt"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv CustomValidator) Validate(i any) error {
	err := cv.validator.Struct(i)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				return fmt.Errorf("%s is required", field)
			case "email":
				return errors.New("invalid email format")
			default:
				return fmt.Errorf("%s is invlaid", field)
			}
		}
	} else {
		return err
	}

	return nil
}

func Init(infra infrastructure.Infrastructure, cfg *config.Configuration) (e *echo.Echo) {
	e = echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Validator = &CustomValidator{validator: validator.New()}

	sqlStore := infra.SQLStore()
	userRepository := user.New(sqlStore.GetDB())
	hasher := hasher.HasherPassword{}
	jwt := jwt.New(cfg)

	registerUsecase := register.New(userRepository, &hasher, jwt)
	registerHandler := registerhandler.New(registerUsecase)

	loginUsecase := login.New(userRepository, &hasher, jwt)
	loginHandler := loginhandler.New(loginUsecase)

	route := e.Group("/v1")
	route.POST("/register", registerHandler.Register)
	route.POST("/login", loginHandler.Login)
	return
}
