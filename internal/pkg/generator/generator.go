package generator

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	defaultLength = 10
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	source      = rand.NewSource(time.Now().UnixNano())
	random      = rand.New(source)
)

func RandSID(n int) string {
	sid := make([]rune, n)
	for i := range sid {
		sid[i] = letterRunes[random.Intn(len(letterRunes))]
	}
	return string(sid)
}

func RandToken() string {
	token := ""
	for i := 1; i < defaultLength; i++ {
		token += strconv.Itoa(random.Intn(9))
	}
	return token
}
