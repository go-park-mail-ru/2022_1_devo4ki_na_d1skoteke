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
	r := mux.NewRouter()

	securityManager := security.NewSimpleSecurityManager()
	userStorage := storage.NewUserCacheStorage(securityManager)
	userService := application.NewUserService(userStorage, securityManager)

	register := interfaces.NewAuthHandler(userService)

	r.HandleFunc("/users/signup", register.SignUp).Methods("POST")

	fmt.Println("Start server at port 3000...")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}
}
