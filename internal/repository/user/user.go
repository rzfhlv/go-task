package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rzfhlv/go-task/internal/model"
)

var (
	createUserQuery = `INSERT INTO users
		(name, email, password)
		VALUES ($1, $2, $3) RETURNING *`

	getByEmailQUery = `SELECT id, name, email, password FROM users WHERE email = $1`
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type User struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) UserRepository {
	return &User{
		db: db,
	}
}

func (u *User) Create(ctx context.Context, user model.User) (model.User, error) {
	result := model.User{}
	err := u.db.Get(&result, createUserQuery, user.Name, user.Email, user.Password)
	if err != nil {
		return model.User{}, err
	}

	return result, nil
}

func (u *User) GetByEmail(ctx context.Context, email string) (model.User, error) {
	result := model.User{}
	err := u.db.Get(&result, getByEmailQUery, email)
	if err != nil {
		return model.User{}, err
	}

	return result, nil
}
