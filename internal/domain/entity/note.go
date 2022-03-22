package entity

type Note struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type Notes struct {
	Notes []Note `json:"notes"`
}
