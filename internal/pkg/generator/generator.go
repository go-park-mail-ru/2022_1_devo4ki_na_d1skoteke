package generator

import "math/rand"

const (
	defaultLength = 10
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func RandSID(n int) string {
	sid := make([]rune, n)
	for i := range sid {
		sid[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(sid)
}

func RandToken() string {
	token := make([]rune, defaultLength)
	for i := range token {
		token[i] = rune(rand.Intn(9))
	}
	return string(token)
}
