package password

import (
	"errors"
	"github.com/dlclark/regexp2"
)

var ErrBadPassword = errors.New("bad password")

const regex = `^(?=.*[0-9])[a-zA-Z0-9!@#$%^&*]{7,30}$`

func ValidatePassword(password string) error {
	regex, err := regexp2.Compile(regex, 0)
	if err != nil {
		return err
	}

	if isMatch, _ := regex.MatchString(password); isMatch {
		return nil
	}
	return ErrBadPassword
}
