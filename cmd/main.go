package main

import (
	"cotion/internal/application/auth"
	"cotion/internal/application/notes"
	"cotion/internal/application/user"
	"cotion/internal/handler"
	"cotion/internal/handler/middleware"
	"cotion/internal/infrastructure/storage"
	"cotion/internal/pkg/security"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	securityManager := security.NewSimpleSecurityManager()

	userStorage := storage.NewUserCacheStorage(securityManager)
	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)
	sessionStorage := storage.NewSessionStorage()

	notesService := notes.NewNotesApp(notesStorage, usersNotesStorage)
	userService := user.NewUserService(userStorage, securityManager)
	authService := auth.NewAuthApp(sessionStorage, userService, securityManager)

	notesHandler := handler.NewNotesHandler(notesService, authService, securityManager)
	registerHandler := handler.NewAuthHandler(userService)
	loginHandler := handler.NewLoginHandler(authService)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", notesHandler.ReceiveSingleNote).Methods("GET")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", notesHandler.UpdateNote).Methods("PUT") //update note data
	routerAPI.HandleFunc("/notes", notesHandler.MainPage).Methods("GET")
	routerAPI.HandleFunc("/note", notesHandler.CreateNote).Methods("POST")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}/delete", notesHandler.DeleteNote).Methods("POST")

	routerAPI.HandleFunc("/users/login", loginHandler.Login).Methods("POST")
	routerAPI.HandleFunc("/users/logout", loginHandler.Logout).Methods("GET")
	routerAPI.HandleFunc("/users/auth", loginHandler.Auth).Methods("GET")
	routerAPI.HandleFunc("/users/signup", registerHandler.SignUp).Methods("POST")

	router.Use(middleware.CorsMiddleware())

	fmt.Println("Start server at port 3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
