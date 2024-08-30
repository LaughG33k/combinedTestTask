package model

type Note struct {
	OwnerUuid string
	Title     string `json:"title"`
	Content   string `json:"content"`
}
