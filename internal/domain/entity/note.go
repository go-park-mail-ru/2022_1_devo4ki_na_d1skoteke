package entity

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	MaxNameLength = 30
	MaxBodyLength = 500
)

var ErrorNoteNameLengthExceedsLimit error = errors.New("note name length exceeds limit")
var ErrorNoteBodyLengthExceedsLimit error = errors.New("note name length exceeds limit")

type Note struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type Notes struct {
	Notes []Note `json:"notes"`
}

type ShortNote struct {
	Name    string `json:"name"`
	Favicon string `json:"favicon"`
	Token   string `json:"token"`
}

type ShortNotes struct {
	ShortNote []ShortNote `json:"notes"`
}

type NoteRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func (n *NoteRequest) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		return err
	}
	if err := n.Validate(); err != nil {
		return err
	}
	return nil
}

func (n *NoteRequest) Validate() error {
	if len(n.Name) > MaxNameLength {
		return ErrorNoteNameLengthExceedsLimit
	}
	if len(n.Name) > MaxBodyLength {
		return ErrorNoteBodyLengthExceedsLimit
	}
	return nil
}
