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

	notesStorage := storage.NewNotesStorage()
	usersNotesStorage := storage.NewUsersNotesStorage(notesStorage)

	notesService := application.NewNotesApp(notesStorage, usersNotesStorage)
	notesHandler := interfaces.NewNotesHandler(notesService)

	userStorage := storage.NewUserCacheStorage(securityManager)
	userService := application.NewUserService(userStorage, securityManager)

	authHandler := interfaces.NewAuthHandler(userService)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", notesHandler.ReceiveSingleNote).Methods("GET")
	routerAPI.HandleFunc("/", notesHandler.MainPage).Methods("GET")

	router.HandleFunc("/users/signup", authHandler.SignUp).Methods("POST")

	fmt.Println("Start server at port 3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
