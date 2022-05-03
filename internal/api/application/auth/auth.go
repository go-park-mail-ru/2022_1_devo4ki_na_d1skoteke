package auth

import (
	"context"
	"cotion/internal/api/application"
	"cotion/internal/api/domain/entity"
	"cotion/internal/pkg/generator"
	"cotion/internal/pkg/security"
	grpcSession "cotion/internal/session/infra/grpc"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	sessionCookie     = "session_id"
	pathSessionCookie = "/api/v1"
	packageName       = "app auth"
)

var ErrNoSession = errors.New("no session")

type AuthApp struct {
	userService     application.UserAppManager
	securityManager security.Manager
	//sessionRepository repository.SessionRepository
	grpcSessionManager grpcSession.AuthCheckerClient
}

func NewAuthApp(grpcManager grpcSession.AuthCheckerClient, userServ application.UserAppManager, secureServ security.Manager) *AuthApp {
	return &AuthApp{
		userService:        userServ,
		securityManager:    secureServ,
		grpcSessionManager: grpcManager,
	}
}

func (au *AuthApp) Login(email string, password string) (*http.Cookie, error) {
	user, err := au.userService.Get(au.securityManager.Hash(email))
	if err != nil {
		return &http.Cookie{}, err
	}

	if err := au.securityManager.ComparePasswords(user.Password, password); err != nil {
		return &http.Cookie{}, err
	}

	SID := generator.RandSID(32)
	ctx := context.Background()
	grpcMessage := &grpcSession.Session{
		SessionID: SID,
		UserID:    user.UserID,
	}
	_, err = au.grpcSessionManager.Create(ctx, grpcMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Login",
		}).Error(err)
		return &http.Cookie{}, err
	}

	cookie := http.Cookie{
		Name:    sessionCookie,
		Value:   SID,
		Expires: time.Now().Add(5 * time.Hour),
		Path:    pathSessionCookie,
	}
	return &cookie, nil
}

func (au *AuthApp) Logout(sessionCookie *http.Cookie) (*http.Cookie, error) {
	ctx := context.Background()
	grpcMessage := &grpcSession.SessionID{
		ID: sessionCookie.Value,
	}
	if _, err := au.grpcSessionManager.Delete(ctx, grpcMessage); err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Logout",
		}).Error(err)
		return &http.Cookie{}, err
	}

	sessionCookie.Expires = time.Now().Add(-time.Hour * 5)
	sessionCookie.Path = pathSessionCookie
	return sessionCookie, nil
}

func (au *AuthApp) Auth(sessionCookie *http.Cookie) (entity.User, bool) {
	ctx := context.Background()
	grpcMessSend := &grpcSession.SessionID{
		ID: sessionCookie.Value,
	}
	grpcMessReceive, err := au.grpcSessionManager.Check(ctx, grpcMessSend)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Auth",
		}).Error(err)
		return entity.User{}, false
	}
	if grpcMessReceive == nil {
		return entity.User{}, false
	}

	user, err := au.userService.Get(grpcMessReceive.UserID)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "Auth",
		}).Error(err)
		return entity.User{}, false
	}

	return user, true
}
