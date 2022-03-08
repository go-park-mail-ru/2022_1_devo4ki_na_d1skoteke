package entity

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) IsEmail() bool {
	return u.Email != ""
}

func (u *User) IsPassword() bool {
	return u.Password != ""
}
