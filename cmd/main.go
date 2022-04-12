package main

import (
	"cotion/internal/application/auth"
	"cotion/internal/application/notes"
	"cotion/internal/application/user"
	"cotion/internal/handler"
	"cotion/internal/handler/middleware"
	"cotion/internal/infrastructure/psql"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/security"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter()

	db, err := psql.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Successful connect to database.")

	securityManager := security.NewSimpleSecurityManager()

	userStorage := psql.NewUserStorage(db)
	notesStorage := psql.NewNotesStorage(db)
	usersNotesStorage := psql.NewUsersNotesStorage(db)
	sessionStorage := storage.NewSessionStorage()

	notesService := notes.NewNotesApp(notesStorage, usersNotesStorage)
	userService := user.NewUserService(userStorage, securityManager)
	authService := auth.NewAuthApp(sessionStorage, userService, securityManager)

	notesHandler := handler.NewNotesHandler(notesService, authService, securityManager)
	userHandler := handler.NewUserHandler(userService)
	loginHandler := handler.NewLoginHandler(authService)

	amw := middleware.NewAuthMiddleware(authService)

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

	router.Use(middleware.CorsMiddleware())

	fmt.Println("Start server at port 3001...")
	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}
}
