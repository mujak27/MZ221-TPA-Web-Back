package model

type Search struct {
	Users []*User `json:"Users"`
	Posts []*Post `json:"Post"`
}
