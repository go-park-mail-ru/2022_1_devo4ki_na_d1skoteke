package main

import (
	"cotion/internal/application/auth"
	"cotion/internal/application/notes"
	"cotion/internal/application/user"
	"cotion/internal/handler"
	"cotion/internal/handler/middleware"
	grpcSession "cotion/internal/infrastructure/grpc"
	"cotion/internal/infrastructure/psql"
	"cotion/internal/infrastructure/s3"
	"cotion/internal/pkg/security"
	"cotion/internal/pkg/xss"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
)

const (
	GRPC_SESSION_URL = "grpcsession"
)

func init() {
	//godotenv.Load(".env_test")
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

	//Minio connect
	imageStorage, err := s3.NewMinioProvider()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successful connect to minio.")

	//Microservice session connect
	grpcSessionConn, err := grpc.Dial(
		os.Getenv(GRPC_SESSION_URL),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error(err)
		log.Fatalf("can't connect to grpc")
	}
	defer grpcSessionConn.Close()
	log.Info("Successful connect to microservice Session.")
	grpcSessManager := grpcSession.NewAuthCheckerClient(grpcSessionConn)

	router := mux.NewRouter()
	securityManager := security.NewSimpleSecurityManager()

	userStorage := psql.NewUserStorage(db)
	notesStorage := psql.NewNotesStorage(db)
	usersNotesStorage := psql.NewUsersNotesStorage(db)
	//sessionStorage := storage.NewSessionStorage()

	notesService := notes.NewNotesApp(notesStorage, usersNotesStorage)
	userService := user.NewUserService(userStorage, imageStorage, securityManager)
	authService := auth.NewAuthApp(grpcSessManager, userService, securityManager)

	notesHandler := handler.NewNotesHandler(notesService, authService, securityManager)
	userHandler := handler.NewUserHandler(userService)
	loginHandler := handler.NewLoginHandler(authService)

	amw := middleware.NewAuthMiddleware(authService)
	xss.NewXssSanitizer()

	routerAPI := router.PathPrefix("/api/v1").Subrouter()

	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.Auth(notesHandler.ReceiveSingleNote)).Methods("GET")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.Auth(notesHandler.UpdateNote)).Methods("PUT") //update note data
	routerAPI.HandleFunc("/notes", amw.Auth(notesHandler.MainPage)).Methods("GET")
	routerAPI.HandleFunc("/note", amw.Auth(notesHandler.CreateNote)).Methods("POST")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.Auth(notesHandler.DeleteNote)).Methods("DELETE")

	routerAPI.HandleFunc("/users/signup", amw.NotAuth(userHandler.SignUp)).Methods("POST")
	routerAPI.HandleFunc("/user", amw.Auth(userHandler.GetUser)).Methods("GET")
	routerAPI.HandleFunc("/user", amw.Auth(userHandler.UpdateUser)).Methods("PUT")
	routerAPI.HandleFunc("/user", amw.Auth(userHandler.DeleteUser)).Methods("DELETE")

	routerAPI.HandleFunc("/users/login", amw.NotAuth(loginHandler.Login)).Methods("POST")
	routerAPI.HandleFunc("/users/logout", amw.Auth(loginHandler.Logout)).Methods("GET")
	routerAPI.HandleFunc("/users/auth", loginHandler.Auth).Methods("GET")

	routerAPI.HandleFunc("/user/avatar", amw.Auth(userHandler.UploadAvatar)).Methods("POST")
	routerAPI.HandleFunc("/user/avatar", amw.Auth(userHandler.DownloadAvatar)).Methods("GET")

	router.Use(middleware.CorsMiddleware())
	//router.Use(middleware.CsrfMiddleware())

	log.Info("Start server at port 3001...")
	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}
