package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordValidator(t *testing.T) {
	cases := map[string]struct {
		password string
		expected error
	}{
		"weak password": {
			password: "a",
			expected: ErrBadPassword,
		},
		"strong password": {
			password: "21ABYHBASD12311213123123asdasd",
			expected: nil,
		},
		"another strong password": {
			password: "Basd213HJJSAD123",
			expected: nil,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			actual := ValidatePassword(tc.password)
			assert.ErrorIs(t, tc.expected, actual)
		})
	}
}
