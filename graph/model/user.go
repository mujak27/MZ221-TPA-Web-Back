package model

var idType = `gorm:"type:varchar(191)`

type User struct {
	ID              string `json:"id" gorm:"type:varchar(191)"`
	Email           string `json:"Email"`
	Password        string `json:"Password"`
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	MidName         string `json:"MidName"`
	IsActive        bool   `json:"IsActive"`
	ProfilePhoto    string `json:"ProfilePhoto"`
	BackgroundPhoto string `json:"BackgroundPhoto"`
	Headline        string `json:"Headline"`
	Pronoun         string `json:"Pronoun"`
	ProfileLink     string `json:"ProfileLink"`
	About           string `json:"About"`
	Location        string `json:"Location"`
}

type Visit struct {
	ID      string `json:"id" gorm:"type:varchar(191)"`
	User1ID string
	User1   *User `json:"User1"`
	User2ID string
	User2   *User `json:"User2"`
}

type Connection struct {
	ID      string `json:"id" gorm:"type:varchar(191)"`
	User1ID string
	User1   *User `json:"User1"`
	User2ID string
	User2   *User `json:"User2"`
}

type Follow struct {
	ID      string `json:"id" gorm:"type:varchar(191)"`
	User1ID string
	User1   *User `json:"User1"`
	User2ID string
	User2   *User `json:"User2"`
}

type InputLogin struct {
	Email    string `json:"email"`
	Passowrd string `json:"passowrd"`
}

type InputRegister struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	MidName   string `json:"MidName"`
}
