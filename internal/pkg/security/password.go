package security

import (
	"bytes"
	"crypto/sha256"
	"errors"
)

type SimpleSecurityManager struct{}

type HashFunc func(string) []byte
type ComparePasswordsFunc func(string, string) error

func NewSimpleSecurityManager() *SimpleSecurityManager {
	return &SimpleSecurityManager{}
}

var wrongPassword = errors.New("wrong password")

func (s SimpleSecurityManager) Hash(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func (s SimpleSecurityManager) ComparePasswords(hashedPassword string, password string) error {
	if bytes.Equal([]byte(hashedPassword), s.Hash(password)) {
		return nil
	}
	return wrongPassword
}

func Hash(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
