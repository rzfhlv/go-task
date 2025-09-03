package auth

import "github.com/rzfhlv/go-task/internal/model"

type AuthResponse struct {
	model.JWT
	User model.User `json:"user"`
}
