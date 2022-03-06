package entity

import (
	"cotion/utils/email"
	"cotion/utils/password"
)

type RegisterUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (u *RegisterUserRequest) Validate() error {
	err := email.ValidateEmail(u.Email)
	if err != nil {
		return err
	}

	if u.ConfirmPassword != u.Password {
		return err
	}

	if err = password.ValidatePassword(u.Password); err != nil {
		return err
	}

	return nil
}
