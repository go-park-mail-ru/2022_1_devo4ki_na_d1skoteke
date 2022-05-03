package main

import (
	"cotion/internal/note/app"
	pb "cotion/internal/note/infra/grpc"
	"cotion/internal/note/infra/psql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

func init() {
	godotenv.Load(".env_test")
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
}

func main() {
	//Postgress connect
	db, err := psql.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Info("Successful connect to database.")

	notesStorage := psql.NewNotesStorage(db)
	usersNotesStorage := psql.NewUsersNotesStorage(db)
	notesServ := app.NewNotesApp(notesStorage, usersNotesStorage)

	//Start server
	lis, err := net.Listen("tcp", ":5002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterNoteServiceServer(server, notesServ)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
