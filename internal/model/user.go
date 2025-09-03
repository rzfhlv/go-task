package model

import "time"

type User struct {
	ID        int64     `json:"id,omitempty" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"password" db:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
