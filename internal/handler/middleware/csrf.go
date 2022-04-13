package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	tokenName = "csrf-token"
)

var secret = "kjkjadjfeqir0w9qrfnkd"

type csrfToken struct {
	value string `json:"csrf-token"`
}

func CsrfMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return enableCSRF(next)
	}
}

func enableCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sCookie, err := r.Cookie(sessionCookie)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			token := r.Header.Get(tokenName)
			if token == "" {
				http.Error(w, "No find csrf token", http.StatusForbidden)
				return
			}
			err := checkToken(sCookie.Value, token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		token, err := createToken(sCookie.Value, time.Now().Add(24*time.Hour).Unix())
		if err != nil {
			log.Error("csrf token creation error:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Add(tokenName, token)
		next.ServeHTTP(w, r)
	})
}

func createToken(userID string, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	data := fmt.Sprintf("%s:%d", userID, tokenExpTime)
	h.Write([]byte(data))
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func checkToken(userID string, inputToken string) error {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return fmt.Errorf("bad token time")
	}

	if tokenExp < time.Now().Unix() {
		return fmt.Errorf("token expired")
	}

	h := hmac.New(sha256.New, []byte(secret))
	data := fmt.Sprintf("%s:%d", userID, tokenExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return fmt.Errorf("cand hex decode token")
	}

	log.Debug(expectedMAC)
	log.Debug(messageMAC)

	if !hmac.Equal(messageMAC, expectedMAC) {
		return fmt.Errorf("invalid token")
	}

	return nil
}
