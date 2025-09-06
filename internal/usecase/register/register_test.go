package register_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/usecase/register"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	cachemocks "github.com/rzfhlv/go-task/internal/repository/cache/mocks"
	usermocks "github.com/rzfhlv/go-task/internal/repository/user/mocks"
	"github.com/rzfhlv/go-task/pkg/errs"
	hashermocks "github.com/rzfhlv/go-task/pkg/hasher/mocks"
	jwtmocks "github.com/rzfhlv/go-task/pkg/jwt/mocks"
)

func TestRegister(t *testing.T) {
	_ = config.Configuration{
		JWT: config.JWTConfiguration{
			ExpiresIn: time.Duration(5 * time.Minute),
		},
	}

	registerRequest := model.Register{
		Name:     "John Doe",
		Email:    "john@mail.com",
		Password: "verysecret",
	}

	passByte, _ := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	hashedPassword := string(passByte)

	userModel := model.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@mail.com",
		Password: hashedPassword,
	}

	jwtModel := model.JWT{
		AccessToken: "token",
		TokenType:   "Bearer",
		ExpiresIn:   300,
	}

	tests := []struct {
		name     string
		request  func(register model.Register) model.Register
		mockDeps func(userRepository *usermocks.MockUserRepository,
			cacheRepository *cachemocks.MockCacheRepository,
			hasher *hashermocks.MockHashPassword,
			jwt *jwtmocks.MockJWTInterface)
		wantUser model.User
		wantJwt  model.JWT
		wantErr  error
	}{
		{
			name: "success",
			request: func(register model.Register) model.Register {
				return register
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository, cacheRepository *cachemocks.MockCacheRepository, hasher *hashermocks.MockHashPassword, jwt *jwtmocks.MockJWTInterface) {
				hasher.On("HashedPassword", mock.MatchedBy(func(password string) bool {
					return password == registerRequest.Password
				})).Return(hashedPassword, nil)

				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == registerRequest.Email
				})).Return(model.User{}, sql.ErrNoRows)

				userRepository.On("Create", context.Background(), mock.MatchedBy(func(register model.Register) bool {
					return register.Name == registerRequest.Name && register.Email == registerRequest.Email && register.Password == hashedPassword
				})).Return(userModel, nil)

				jwt.On("Generate", userModel, mock.Anything).Return(jwtModel, nil)

				cacheRepository.On("Set", context.Background(), mock.Anything, mock.MatchedBy(func(userId int64) bool {
					return userId == userModel.ID
				}), mock.Anything).Return(nil)
			},
			wantUser: userModel,
			wantJwt:  jwtModel,
			wantErr:  nil,
		},
		{
			name: "error when set cache for token",
			request: func(register model.Register) model.Register {
				return register
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository, cacheRepository *cachemocks.MockCacheRepository, hasher *hashermocks.MockHashPassword, jwt *jwtmocks.MockJWTInterface) {
				hasher.On("HashedPassword", mock.MatchedBy(func(password string) bool {
					return password == registerRequest.Password
				})).Return(hashedPassword, nil)

				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == registerRequest.Email
				})).Return(model.User{}, sql.ErrNoRows)

				userRepository.On("Create", context.Background(), mock.MatchedBy(func(register model.Register) bool {
					return register.Name == registerRequest.Name && register.Email == registerRequest.Email && register.Password == hashedPassword
				})).Return(userModel, nil)

				jwt.On("Generate", userModel, mock.Anything).Return(jwtModel, nil)

				cacheRepository.On("Set", context.Background(), mock.Anything, mock.MatchedBy(func(userId int64) bool {
					return userId == userModel.ID
				}), mock.Anything).Return(errors.New("some error"))
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when generate jwt token",
			request: func(register model.Register) model.Register {
				return register
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository, cacheRepository *cachemocks.MockCacheRepository, hasher *hashermocks.MockHashPassword, jwt *jwtmocks.MockJWTInterface) {
				hasher.On("HashedPassword", mock.MatchedBy(func(password string) bool {
					return password == registerRequest.Password
				})).Return(hashedPassword, nil)

				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == registerRequest.Email
				})).Return(model.User{}, sql.ErrNoRows)

				userRepository.On("Create", context.Background(), mock.MatchedBy(func(register model.Register) bool {
					return register.Name == registerRequest.Name && register.Email == registerRequest.Email && register.Password == hashedPassword
				})).Return(userModel, nil)

				jwt.On("Generate", userModel, mock.Anything).Return(model.JWT{}, errors.New("some error"))

				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusUnprocessableEntity, "failed generated token"),
		},
		{
			name: "error when create user",
			request: func(register model.Register) model.Register {
				return register
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository, cacheRepository *cachemocks.MockCacheRepository, hasher *hashermocks.MockHashPassword, jwt *jwtmocks.MockJWTInterface) {
				hasher.On("HashedPassword", mock.MatchedBy(func(password string) bool {
					return password == registerRequest.Password
				})).Return(hashedPassword, nil)

				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == registerRequest.Email
				})).Return(model.User{}, sql.ErrNoRows)

				userRepository.On("Create", context.Background(), mock.MatchedBy(func(register model.Register) bool {
					return register.Name == registerRequest.Name && register.Email == registerRequest.Email && register.Password == hashedPassword
				})).Return(model.User{}, errors.New("some error"))

				jwt.AssertNotCalled(t, "Generate")
				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when user already exists",
			request: func(register model.Register) model.Register {
				return register
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository, cacheRepository *cachemocks.MockCacheRepository, hasher *hashermocks.MockHashPassword, jwt *jwtmocks.MockJWTInterface) {
				hasher.On("HashedPassword", mock.MatchedBy(func(password string) bool {
					return password == registerRequest.Password
				})).Return(hashedPassword, nil)

				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == registerRequest.Email
				})).Return(userModel, nil)

				userRepository.AssertNotCalled(t, "Create")
				jwt.AssertNotCalled(t, "Generate")
				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusUnprocessableEntity, "email already exists"),
		},
		{
			name: "error when hashed password",
			request: func(register model.Register) model.Register {
				return register
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository, cacheRepository *cachemocks.MockCacheRepository, hasher *hashermocks.MockHashPassword, jwt *jwtmocks.MockJWTInterface) {
				hasher.On("HashedPassword", mock.MatchedBy(func(password string) bool {
					return password == registerRequest.Password
				})).Return("", errors.New("some error"))

				userRepository.AssertNotCalled(t, "GetByEmail")
				userRepository.AssertNotCalled(t, "Create")
				jwt.AssertNotCalled(t, "Generate")
				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusUnprocessableEntity, "hasher error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.All("../../../")
			userRepository := usermocks.MockUserRepository{}
			cacheRepository := cachemocks.MockCacheRepository{}
			hasher := hashermocks.MockHashPassword{}
			jwt := jwtmocks.MockJWTInterface{}

			request := tt.request(registerRequest)

			tt.mockDeps(&userRepository, &cacheRepository, &hasher, &jwt)

			usecase := register.New(&userRepository, &cacheRepository, &hasher, &jwt)
			userResult, jwtResult, err := usecase.Register(context.Background(), request)

			assert.Equal(t, tt.wantUser, userResult)
			assert.Equal(t, tt.wantJwt, jwtResult)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
