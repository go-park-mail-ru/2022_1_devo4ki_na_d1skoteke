package auth

import (
	"cotion/internal/application"
	"cotion/internal/domain/entity"
	"cotion/internal/domain/repository"
	"cotion/internal/pkg/generator"
	"cotion/internal/pkg/security"
	"errors"
	"net/http"
	"time"
)

const (
	sessionCookie     = "session_id"
	pathSessionCookie = "/api/v1"
)

type AuthApp struct {
	userService       application.UserAppManager
	securityManager   security.Manager
	sessionRepository repository.SessionRepository
}

func NewAuthApp(sessionRepo repository.SessionRepository, userServ application.UserAppManager, secureServ security.Manager) *AuthApp {
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
		Name:    sessionCookie,
		Value:   session.SID,
		Expires: time.Now().Add(5 * time.Hour),
		Path:    pathSessionCookie,
	}
	return &cookie, nil
}

func (au *AuthApp) Logout(sessionCookie *http.Cookie) (*http.Cookie, error) {
	if _, ok := au.sessionRepository.HasSession(sessionCookie.Value); !ok {
		return &http.Cookie{}, errors.New("no session")
	}

	au.sessionRepository.DeleteSession(sessionCookie.Value)
	sessionCookie.Expires = time.Now().Add(-time.Hour * 5)
	sessionCookie.Path = pathSessionCookie
	return sessionCookie, nil
}

func (au *AuthApp) Auth(sessionCookie *http.Cookie) (entity.User, bool) {
	session, ok := au.sessionRepository.HasSession(sessionCookie.Value)
	if !ok {
		return entity.User{}, false
	}

	user, err := au.userService.GetUser(session.UserEmail)
	if err != nil {
		return entity.User{}, false
	}

	return user, true
}