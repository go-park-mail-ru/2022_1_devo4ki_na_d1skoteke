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

	amw := middleware.NewAuthMiddleware(authService)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.AuthMiddleware(true, notesHandler.ReceiveSingleNote)).Methods("GET")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", amw.AuthMiddleware(true, notesHandler.UpdateNote)).Methods("PUT") //update note data
	routerAPI.HandleFunc("/notes", amw.AuthMiddleware(true, notesHandler.MainPage)).Methods("GET")
	routerAPI.HandleFunc("/note", amw.AuthMiddleware(true, notesHandler.CreateNote)).Methods("POST")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}/delete", amw.AuthMiddleware(true, notesHandler.DeleteNote)).Methods("POST")

	routerAPI.HandleFunc("/users/login", amw.AuthMiddleware(false, loginHandler.Login)).Methods("POST")
	routerAPI.HandleFunc("/users/logout", amw.AuthMiddleware(true, loginHandler.Logout)).Methods("GET")
	routerAPI.HandleFunc("/users/auth", loginHandler.Auth).Methods("GET")
	routerAPI.HandleFunc("/users/signup", amw.AuthMiddleware(false, registerHandler.SignUp)).Methods("POST")

	router.Use(middleware.CorsMiddleware())

	fmt.Println("Start server at port 3001...")
	err := http.ListenAndServe(":3001", router)
	if err != nil {
		log.Fatal(err)
	}
}
