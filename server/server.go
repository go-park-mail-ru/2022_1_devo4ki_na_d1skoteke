package main

import (
	"cotion/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.HomeHandler)
	router.HandleFunc("/api/v1/users/signup", handlers.SignupAPI).Methods("POST")
	router.HandleFunc("/api/v1/users/login", handlers.LoginAPI).Methods("POST")
	router.HandleFunc("/api/v1/users/logout", handlers.LogoutAPI)
	router.HandleFunc("/api/v1/notes", handlers.NotesAPI).Methods("GET")
	router.HandleFunc("/api/v1/note/{note-token:[0-9]+}", handlers.SingleNoteAPI).Methods("GET")

	fmt.Println("Start server at port 3000...")
	http.ListenAndServe(":3000", router)
}
