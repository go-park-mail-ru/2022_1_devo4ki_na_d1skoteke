package main

import (
	"cotion/internal/application/auth"
	"cotion/internal/application/notes"
	userapp "cotion/internal/application/user"
	"cotion/internal/handler"
	"cotion/internal/handler/middleware"
	"cotion/internal/infrastructure/psql"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/security"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "0.0.0.0"
	port     = 5432
	user     = "user"
	password = "love"
	dbname   = "cotion"
)

func main() {
	router := mux.NewRouter()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful connect to database.")

	securityManager := security.NewSimpleSecurityManager()

	userStorage := storage.NewUserCacheStorage(securityManager)
	notesStorage := psql.NewNotesStorage(db)
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	sessionStorage := storage.NewSessionStorage()

	notesService := notes.NewNotesApp(notesStorage, usersNotesStorage)
	userService := userapp.NewUserService(userStorage, securityManager)
	authService := auth.NewAuthApp(sessionStorage, userService, securityManager)

	notesHandler := handler.NewNotesHandler(notesService, authService, securityManager)
	registerHandler := handler.NewAuthHandler(userService)
	loginHandler := handler.NewLoginHandler(authService)

	amw := middleware.NewAuthMiddleware(authService)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.Auth(notesHandler.ReceiveSingleNote)).Methods("GET")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.Auth(notesHandler.UpdateNote)).Methods("PUT") //update note data
	routerAPI.HandleFunc("/notes", amw.Auth(notesHandler.MainPage)).Methods("GET")
	routerAPI.HandleFunc("/note", amw.Auth(notesHandler.CreateNote)).Methods("POST")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.Auth(notesHandler.DeleteNote)).Methods("DELETE")

	routerAPI.HandleFunc("/users/login", amw.NotAuth(loginHandler.Login)).Methods("POST")
	routerAPI.HandleFunc("/users/logout", amw.Auth(loginHandler.Logout)).Methods("GET")
	routerAPI.HandleFunc("/users/auth", loginHandler.Auth).Methods("GET")
	routerAPI.HandleFunc("/users/signup", amw.NotAuth(registerHandler.SignUp)).Methods("POST")

	router.Use(middleware.CorsMiddleware())

	fmt.Println("Start server at port 3001...")
	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}
