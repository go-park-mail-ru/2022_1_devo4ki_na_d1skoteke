package application

import (
	"cotion/domain/repository"
	"cotion/infrastructure/security"
	"cotion/utils/generator"
	"errors"
	"net/http"
	"time"
)

const (
	cookieSession = "session_id"
)

type AuthApp struct {
	userService       UserAppManager
	securityManager   security.Manager
	sessionRepository repository.SessionRepository
}

func NewAuthApp(sessionRepo repository.SessionRepository, userServ UserAppManager, secureServ security.Manager) *AuthApp {
	return &AuthApp{
		userService:       userServ,
		securityManager:   secureServ,
		sessionRepository: sessionRepo,
	}
}

func (au *AuthApp) Login(email string, password string) (*http.Cookie, error) {
	user, err := au.userService.GetUser(email)
	if err != nil {
		return &http.Cookie{}, err
	}

	if err = au.securityManager.ComparePasswords(user.Password, password); err != nil {
		return &http.Cookie{}, err
	}

	SID := generator.RandSID(32)
	session, err := au.sessionRepository.NewSession(SID, user)
	if err != nil {
		return &http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:    cookieSession,
		Value:   session.SID,
		Expires: time.Now().Add(10 * time.Hour),
	}
	return &cookie, nil
}

func (au *AuthApp) Logout(sessionCookie *http.Cookie) (*http.Cookie, error) {
	if _, ok := au.sessionRepository.HasSession(sessionCookie.Value); !ok {
		return &http.Cookie{}, errors.New("no session")
	}

	au.sessionRepository.DeleteSession(sessionCookie.Value)
	sessionCookie.Expires = time.Now().AddDate(0, 0, -1)
	return sessionCookie, nil
}

func (au *AuthApp) Auth(SID string) bool {
	_, ok := au.sessionRepository.HasSession(SID)
	return ok
}
