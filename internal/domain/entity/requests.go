package entity

import (
	"cotion/internal/pkg/email"
	"cotion/internal/pkg/password"
	"errors"
)

type RegisterUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

var DifferentPasswords = errors.New("different passwords")

func (u *RegisterUserRequest) Validate() error {
	err := email.ValidateEmail(u.Email)
	if err != nil {
		return err
	}

	if u.ConfirmPassword != u.Password {
		return DifferentPasswords
	}

	if err = password.ValidatePassword(u.Password); err != nil {
		return err
	}

	return nil
}
