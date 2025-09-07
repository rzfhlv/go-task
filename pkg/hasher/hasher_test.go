package hasher_test

import (
	"errors"
	"testing"

	"github.com/rzfhlv/go-task/pkg/hasher"
	"github.com/stretchr/testify/assert"
)

func TestHasherHashedPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "success",
			password: "password",
			wantErr:  nil,
		},
		{
			name:     "error hashed password",
			password: "passwordghhghgkgkggkjgjkgjkgkjgkgjkgggkjgkgjkgkjjkjljlkjkljljkljlkjljljklj",
			wantErr:  errors.New("bcrypt: password length exceeds 72 bytes"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := hasher.HasherPassword{}
			hashed, err := hasher.HashedPassword(tt.password)

			assert.NotNil(t, hashed)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestHasherVerifyPassword(t *testing.T) {
	tests := []struct {
		name, hashed, password string
		wantErr                error
	}{
		{
			name:     "success",
			hashed:   "$2a$10$d3.zWWlz0tAnXis7fAJulumr2JHT5YDoZ7OzY9yJcx1TmQhS7c4mO",
			password: "password",
			wantErr:  nil,
		},
		{
			name:     "error when verify password",
			hashed:   "$2a$10$d3.zWWlz0tAnXis7fAJulumr2JHT5YDoZ7OzY9yJcx1TmQhS7c4mO",
			password: "invalidpassword",
			wantErr:  errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := hasher.HasherPassword{}
			err := hasher.VerifyPassword(tt.hashed, tt.password)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
