package entity

import "errors"

type Note struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type NoteShort struct {
	Name    string `json:"name"`
	Favicon string `json:"favicon"`
	Token   string `json:"token"`
}

type Notes struct {
	Notes []Note `json:"notes"`
}

const (
	MaxNameLength = 30
	MaxBodyLength = 500
)

var ErrorNoteNameLengthExceedsLimit error = errors.New("note name length exceeds limit")
var ErrorNoteBodyLengthExceedsLimit error = errors.New("note name length exceeds limit")

func (n *Note) Validate() error {
	if len(n.Name) > MaxNameLength {
		return ErrorNoteNameLengthExceedsLimit
	}
	if len(n.Name) > MaxBodyLength {
		return ErrorNoteBodyLengthExceedsLimit
	}
	return nil
}
