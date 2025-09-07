package user_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/rzfhlv/go-task/internal/model"
	"github.com/rzfhlv/go-task/internal/repository/user"
	"github.com/stretchr/testify/assert"
)

var (
	now = time.Date(2023, time.August, 15, 12, 0, 0, 0, time.UTC)

	userModel = model.User{
		ID:        1,
		Name:      "John",
		Email:     "john@mail.com",
		Password:  "verysecret",
		CreatedAt: now,
	}
)

func TestUserCreate(t *testing.T) {
	register := model.Register{
		Name:     "John",
		Email:    "john@mail.com",
		Password: "verysecret",
	}

	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult model.User
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at"}).
					AddRow(userModel.ID, userModel.Name, userModel.Email, userModel.Password, userModel.CreatedAt)

				s.ExpectQuery(`INSERT INTO users 
				(name, email, password) 
				VALUES ($1, $2, $3) RETURNING *`).
					WithArgs(userModel.Name, userModel.Email, userModel.Password).
					WillReturnRows(rows)
			},
			wantResult: userModel,
			wantErr:    nil,
		},
		{
			name: "error when insert",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(`INSERT INTO users 
				(name, email, password) 
				VALUES ($1, $2, $3) RETURNING *`).
					WithArgs(userModel.Name, userModel.Email, userModel.Password).
					WillReturnError(sql.ErrConnDone)
			},
			wantResult: model.User{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := user.New(db)
			result, err := r.Create(context.Background(), register)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}

func TestUserGetByEmail(t *testing.T) {
	email := "john@mail.com"

	tests := []struct {
		name       string
		beforeTest func(s sqlmock.Sqlmock)
		wantResult model.User
		wantErr    error
	}{
		{
			name: "success",
			beforeTest: func(s sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at"}).
					AddRow(userModel.ID, userModel.Name, userModel.Email, userModel.Password, userModel.CreatedAt)

				s.ExpectQuery("SELECT id, name, email, password, created_at FROM users WHERE email = $1").
					WithArgs(email).WillReturnRows(rows)
			},
			wantResult: userModel,
			wantErr:    nil,
		},
		{
			name: "error when get by email",
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectQuery("SELECT id, name, email, password, created_at FROM users WHERE email = $1").
					WithArgs(email).WillReturnError(sql.ErrConnDone)
			},
			wantResult: model.User{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer mockDB.Close()

			db := sqlx.NewDb(mockDB, "sqlmock")

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			r := user.New(db)
			result, err := r.GetByEmail(context.Background(), email)

			assert.Equal(t, tt.wantResult, result)
			assert.Equal(t, tt.wantErr, err)

			if err := mockSQL.ExpectationsWereMet(); err != nil {
				assert.Error(t, err)
			}
		})
	}
}
