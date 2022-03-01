package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func SignupAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "signup")
}

func LoginAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login")
}

func LogoutAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "logout")
}

func NotesAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "notes")
}

func SingleNoteAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "note: "+r.RequestURI)
}
