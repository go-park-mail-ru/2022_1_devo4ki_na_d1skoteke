package password

import (
	passValidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 50.

func ValidatePassword(password string) error {
	return passValidator.Validate(password, minEntropyBits)
}
