package model

import "time"

type Task struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description" db:"description" validate:"required"`
	Status      string    `json:"status" db:"status"`
	UserID      int64     `json:"-" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
