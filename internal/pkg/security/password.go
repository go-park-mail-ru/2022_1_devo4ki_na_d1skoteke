package security

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

type SimpleSecurityManager struct{}

type HashFunc func(string) []byte
type ComparePasswordsFunc func(string, string) error

func NewSimpleSecurityManager() *SimpleSecurityManager {
	return &SimpleSecurityManager{}
}

var wrongPassword = errors.New("wrong password")

func (s SimpleSecurityManager) Hash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (s SimpleSecurityManager) ComparePasswords(hashedPassword string, password string) error {
	if strings.Compare(hashedPassword, s.Hash(password)) == 0 {
		return nil
	}
	return wrongPassword
}

func Hash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
