package entity

import (
	"cotion/internal/pkg/email"
	"cotion/internal/pkg/password"
	"encoding/json"
	"errors"
	"net/http"
)

type User struct {
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

func (u *User) IsEmail() bool {
	return u.Email != ""
}

func (u *User) IsPassword() bool {
	return u.Password != ""
}

type UserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

var ErrDiffPasswords = errors.New("passwords don't match")

func (u *UserRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return err
	}

	return u.Validate()
}

func (u *UserRequest) Validate() error {
	err := email.ValidateEmail(u.Email)
	if err != nil {
		return err
	}

	if u.ConfirmPassword != u.Password {
		return ErrDiffPasswords
	}

	if err = password.ValidatePassword(u.Password); err != nil {
		return err
	}

	return nil
}
