package generator

import "math/rand"

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
