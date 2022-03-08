package main

import (
	"cotion/application"
	"cotion/infrastructure/security"
	storage "cotion/infrastructure/storage"
	"cotion/interfaces"
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
  
	notesService := application.NewNotesApp(notesStorage, usersNotesStorage)
	userService := application.NewUserService(userStorage, securityManager)
	authService := application.NewAuthApp(sessionStorage, userService, securityManager)
	
  notesHandler := interfaces.NewNotesHandler(notesService)
	registerHandler := interfaces.NewAuthHandler(userService)
 	loginHandler := interfaces.NewLoginHandler(authService)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", notesHandler.ReceiveSingleNote).Methods("GET")
	routerAPI.HandleFunc("/notes", notesHandler.MainPage).Methods("GET")
  routerAPI.HandleFunc("/users/login", loginHandler.Login).Methods("POST")
	routerAPI.HandleFunc("/users/logout", loginHandler.Logout).Methods("GET")
	routerAPI.HandleFunc("/users/auth", loginHandler.Auth).Methods("GET")
	routerAPI.HandleFunc("/users/signup", registerHandler.SignUp).Methods("POST")

	fmt.Println("Start server at port 3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
