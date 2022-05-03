package storage

import (
	"context"
	pb "cotion/internal/session/infra/grpc"
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

var ErrCreateSession = errors.New("already has session with this SID")

type SessionManager struct {
	pb.UnimplementedAuthCheckerServer
	data sync.Map
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (sm *SessionManager) Create(ctx context.Context, in *pb.Session) (*pb.Nothing, error) {
	_, loaded := sm.data.LoadOrStore(in.SessionID, in)
	if loaded {
		log.WithFields(log.Fields{
			"package":  "storage session",
			"function": "NewSession",
			"sid":      in.SessionID,
			"userID":   in.UserID,
		}).Error(ErrCreateSession)
		return &pb.Nothing{Status: false}, ErrCreateSession
	}
	return &pb.Nothing{Status: true}, nil
}

func (sm *SessionManager) Check(ctx context.Context, in *pb.SessionID) (*pb.Session, error) {
	session, ok := sm.data.Load(in.ID)
	if !ok {
		return nil, nil
	}
	return session.(*pb.Session), nil
}

func (sm *SessionManager) Delete(ctx context.Context, in *pb.SessionID) (*pb.Nothing, error) {
	sm.data.LoadAndDelete(in.ID)
	return &pb.Nothing{Status: true}, nil
}
