package rest

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rzfhlv/go-task/config"
	loginhandler "github.com/rzfhlv/go-task/internal/handler/login"
	logouthandler "github.com/rzfhlv/go-task/internal/handler/logout"
	registerhandler "github.com/rzfhlv/go-task/internal/handler/register"
	taskhandler "github.com/rzfhlv/go-task/internal/handler/task"
	"github.com/rzfhlv/go-task/internal/infrastructure"
	"github.com/rzfhlv/go-task/internal/repository/cache"
	"github.com/rzfhlv/go-task/internal/repository/task"
	"github.com/rzfhlv/go-task/internal/repository/user"
	"github.com/rzfhlv/go-task/internal/usecase/login"
	"github.com/rzfhlv/go-task/internal/usecase/logout"
	"github.com/rzfhlv/go-task/internal/usecase/register"
	taskusecase "github.com/rzfhlv/go-task/internal/usecase/task"
	"github.com/rzfhlv/go-task/pkg/hasher"
	"github.com/rzfhlv/go-task/pkg/jwt"
	"github.com/rzfhlv/go-task/pkg/middleware/auth"
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
	memStore := infra.MemStore()
	userRepository := user.New(sqlStore.GetDB())
	cacheRepository := cache.New(memStore.GetClient())
	taskRepository := task.New(sqlStore.GetDB())

	hasher := hasher.HasherPassword{}
	jwt := jwt.New(cfg)

	middleware := auth.New(cacheRepository, jwt)

	registerUsecase := register.New(userRepository, cacheRepository, &hasher, jwt)
	registerHandler := registerhandler.New(registerUsecase)

	loginUsecase := login.New(userRepository, cacheRepository, &hasher, jwt)
	loginHandler := loginhandler.New(loginUsecase)

	logoutUsecase := logout.New(cacheRepository)
	logoutHandler := logouthandler.New(logoutUsecase)

	taskUsecase := taskusecase.New(taskRepository)
	taskHandler := taskhandler.New(taskUsecase)

	route := e.Group("/v1")
	route.POST("/register", registerHandler.Register)
	route.POST("/login", loginHandler.Login)
	route.POST("/logout", logoutHandler.Logout, middleware.Bearer)

	task := route.Group("/tasks", middleware.Bearer)
	task.POST("", taskHandler.Create)
	task.GET("", taskHandler.GetByUserID)
	task.GET("/:id", taskHandler.GetByID)
	task.PUT("/:id", taskHandler.Update)
	task.DELETE("/:id", taskHandler.Delete)

	return
}
