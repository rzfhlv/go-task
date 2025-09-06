package login_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/rzfhlv/go-task/config"
	"github.com/rzfhlv/go-task/internal/model"
	cachemocks "github.com/rzfhlv/go-task/internal/repository/cache/mocks"
	usermocks "github.com/rzfhlv/go-task/internal/repository/user/mocks"
	"github.com/rzfhlv/go-task/internal/usecase/login"
	"github.com/rzfhlv/go-task/pkg/errs"
	hashermocks "github.com/rzfhlv/go-task/pkg/hasher/mocks"
	jwtmocks "github.com/rzfhlv/go-task/pkg/jwt/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	_ = config.Configuration{
		JWT: config.JWTConfiguration{
			ExpiresIn: time.Duration(5 * time.Minute),
		},
	}

	loginRequest := model.Login{
		Email:    "john@mail.com",
		Password: "verysecret",
	}

	passByte, _ := bcrypt.GenerateFromPassword([]byte(loginRequest.Password), bcrypt.DefaultCost)
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
		request  func(login model.Login) model.Login
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
			request: func(login model.Login) model.Login {
				return login
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository,
				cacheRepository *cachemocks.MockCacheRepository,
				hasher *hashermocks.MockHashPassword,
				jwt *jwtmocks.MockJWTInterface) {
				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == loginRequest.Email
				})).
					Return(userModel, nil)

				hasher.On("VerifyPassword", mock.MatchedBy(func(hashedPassword string) bool {
					return hashedPassword == userModel.Password
				}), mock.MatchedBy(func(password string) bool {
					return password == loginRequest.Password
				})).Return(nil)

				jwt.On("Generate", userModel, mock.Anything).Return(jwtModel, nil)

				cacheRepository.On("Set", mock.Anything, mock.Anything, mock.MatchedBy(func(userId int64) bool {
					return userId == userModel.ID
				}), mock.Anything).Return(nil)
			},
			wantUser: userModel,
			wantJwt:  jwtModel,
			wantErr:  nil,
		},
		{
			name: "error when set cache for token",
			request: func(login model.Login) model.Login {
				return login
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository,
				cacheRepository *cachemocks.MockCacheRepository,
				hasher *hashermocks.MockHashPassword,
				jwt *jwtmocks.MockJWTInterface) {
				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == loginRequest.Email
				})).
					Return(userModel, nil)

				hasher.On("VerifyPassword", mock.MatchedBy(func(hashedPassword string) bool {
					return hashedPassword == userModel.Password
				}), mock.MatchedBy(func(password string) bool {
					return password == loginRequest.Password
				})).Return(nil)

				jwt.On("Generate", userModel, mock.Anything).Return(jwtModel, nil)

				cacheRepository.On("Set", mock.Anything, mock.Anything, mock.MatchedBy(func(userId int64) bool {
					return userId == userModel.ID
				}), mock.Anything).Return(errors.New("some error"))
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when generate jwt token",
			request: func(login model.Login) model.Login {
				return login
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository,
				cacheRepository *cachemocks.MockCacheRepository,
				hasher *hashermocks.MockHashPassword,
				jwt *jwtmocks.MockJWTInterface) {
				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == loginRequest.Email
				})).
					Return(userModel, nil)

				hasher.On("VerifyPassword", mock.MatchedBy(func(hashedPassword string) bool {
					return hashedPassword == userModel.Password
				}), mock.MatchedBy(func(password string) bool {
					return password == loginRequest.Password
				})).Return(nil)

				jwt.On("Generate", userModel, mock.Anything).Return(model.JWT{}, errors.New("some error"))

				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusUnauthorized, "invalid credentials"),
		},
		{
			name: "error when verify password",
			request: func(login model.Login) model.Login {
				return login
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository,
				cacheRepository *cachemocks.MockCacheRepository,
				hasher *hashermocks.MockHashPassword,
				jwt *jwtmocks.MockJWTInterface) {
				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == loginRequest.Email
				})).
					Return(userModel, nil)

				hasher.On("VerifyPassword", mock.MatchedBy(func(hashedPassword string) bool {
					return hashedPassword == userModel.Password
				}), mock.MatchedBy(func(password string) bool {
					return password == loginRequest.Password
				})).Return(bcrypt.ErrMismatchedHashAndPassword)

				jwt.AssertNotCalled(t, "Generate")
				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusUnauthorized, "invalid credentials"),
		},
		{
			name: "error when get user by email",
			request: func(login model.Login) model.Login {
				return login
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository,
				cacheRepository *cachemocks.MockCacheRepository,
				hasher *hashermocks.MockHashPassword,
				jwt *jwtmocks.MockJWTInterface) {
				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == loginRequest.Email
				})).
					Return(model.User{}, errors.New("some error"))

				hasher.AssertNotCalled(t, "VerifyPassword")
				jwt.AssertNotCalled(t, "Generate")
				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusInternalServerError, "something went wrong"),
		},
		{
			name: "error when get user by email and sql no result message",
			request: func(login model.Login) model.Login {
				return login
			},
			mockDeps: func(userRepository *usermocks.MockUserRepository,
				cacheRepository *cachemocks.MockCacheRepository,
				hasher *hashermocks.MockHashPassword,
				jwt *jwtmocks.MockJWTInterface) {
				userRepository.On("GetByEmail", context.Background(), mock.MatchedBy(func(email string) bool {
					return email == loginRequest.Email
				})).
					Return(model.User{}, sql.ErrNoRows)

				hasher.AssertNotCalled(t, "VerifyPassword")
				jwt.AssertNotCalled(t, "Generate")
				cacheRepository.AssertNotCalled(t, "Set")
			},
			wantUser: model.User{},
			wantJwt:  model.JWT{},
			wantErr:  errs.NewErrs(http.StatusUnauthorized, "unauthorized"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.All("../../../")
			userRepository := usermocks.MockUserRepository{}
			cacheRepository := cachemocks.MockCacheRepository{}
			hasher := hashermocks.MockHashPassword{}
			jwt := jwtmocks.MockJWTInterface{}

			request := tt.request(loginRequest)

			tt.mockDeps(&userRepository, &cacheRepository, &hasher, &jwt)

			usecase := login.New(&userRepository, &cacheRepository, &hasher, &jwt)
			userResult, jwtResult, err := usecase.Login(context.Background(), request)

			assert.Equal(t, tt.wantUser, userResult)
			assert.Equal(t, tt.wantJwt, jwtResult)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
