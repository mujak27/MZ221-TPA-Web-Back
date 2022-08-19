// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Post struct {
	ID     string `gorm:"type:varchar(191)"`
	Text    string  `json:"Text"`
	SenderId string 
	Sender *User  `json:"Sender" gorm:"reference:User"`
}


// type Post struct {
// 	ID     string `gorm:"type:varchar(191)"`
// 	Text    string  `json:"Text"`
// 	Senders []*User `json:"Senders"   gorm:"many2many:post_senders;"`
// }

// type PostSender struct {
// 	ID     string `gorm:"type:varchar(191)"`
// 	UserId string `json:"UserId" gorm:"reference:User;type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// 	PostId string `json:"PostId" gorm:"reference:Post;type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
// }



type InputPost struct {
	Text   string `json:"text"`
	Offset int    `json:"Offset"`
	Limit  int    `json:"Limit"`
}
