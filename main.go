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

	notesStorage := storage.NewNotesStorage()
	notesService := application.NewNotesApp(notesStorage)
	notesHandler := interfaces.NewNotesHandler(notesService)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", notesHandler.ReceiveSingleNote).Methods("GET")

	securityManager := security.NewSimpleSecurityManager()
	userStorage := storage.NewUserCacheStorage(securityManager)
	userService := application.NewUserService(userStorage, securityManager)
	register := interfaces.NewAuthHandler(userService)

	router.HandleFunc("/users/signup", register.SignUp).Methods("POST")

	fmt.Println("Start server at port 3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
