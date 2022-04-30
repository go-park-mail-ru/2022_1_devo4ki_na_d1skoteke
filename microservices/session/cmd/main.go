package main

import (
	mssession "cotion/microservices/session"
	pb "cotion/microservices/session/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterAuthCheckerServer(server, mssession.NewSessionManager())
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
