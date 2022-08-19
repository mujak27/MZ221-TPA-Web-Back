package model

type Search struct {
	Users []*User `json:"Users"`
	Post  []*Post `json:"Post"`
}
