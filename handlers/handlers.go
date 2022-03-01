package handlers

import (
	"fmt"
	"net/http"
)

type HandlerAPI struct {
}

func CreateHandlerAPI() *HandlerAPI {
	return new(HandlerAPI)
}

func (h *HandlerAPI) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}

func (h *HandlerAPI) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "signup")
}

func (h *HandlerAPI) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login")
}

func (h *HandlerAPI) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "logout")
}

func (h *HandlerAPI) Notes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "notes")
}

func (h *HandlerAPI) SingleNote(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "note: "+r.RequestURI)
}
