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
	MAX_NAME_LENGTH = 30
	MAX_BODY_LENGTH = 500
)

func (n *Note) Validate() error {
	if len(n.Name) > MAX_NAME_LENGTH {
		return errors.New("note name length exceeds limit")
	}
	if len(n.Name) > MAX_BODY_LENGTH {
		return errors.New("note name length exceeds limit")
	}
	return nil
}
