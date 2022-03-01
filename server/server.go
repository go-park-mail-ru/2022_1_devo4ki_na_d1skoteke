package main

import (
	"cotion/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	var handler handlers.HandlerAPI
	router.HandleFunc("/", handler.HomeHandler)

	routerAPI := router.PathPrefix("/api/v1").Subrouter()
	routerAPI.HandleFunc("/users/signup", handler.Signup).Methods("POST")
	routerAPI.HandleFunc("/users/login", handler.Login).Methods("POST")
	routerAPI.HandleFunc("/users/logout", handler.Logout)
	routerAPI.HandleFunc("/notes", handler.Notes).Methods("GET")
	routerAPI.HandleFunc("/note/{note-token:[0-9]+}", handler.SingleNote).Methods("GET")

	fmt.Println("Start server at port 3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
