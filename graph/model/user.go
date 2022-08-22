package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

var idType = `gorm:"type:varchar(191)`

type User struct {
	ID              string        `json:"id" gorm:"type:varchar(191)"`
	Email           string        `json:"Email"`
	Password        string        `json:"Password"`
	FirstName       string        `json:"FirstName"`
	LastName        string        `json:"LastName"`
	MidName         string        `json:"MidName"`
	IsActive        bool          `json:"IsActive"`
	ProfilePhoto    string        `json:"ProfilePhoto"`
	BackgroundPhoto string        `json:"BackgroundPhoto"`
	Headline        string        `json:"Headline"`
	Pronoun         string        `json:"Pronoun"`
	ProfileLink     string        `json:"ProfileLink"`
	About           string        `json:"About"`
	Location        string        `json:"Location"`
	Visits          []*User       `json:"Visit" gorm:"many2many:user_visits"`
	Follows         []*User       `json:"Follow" gorm:"many2many:user_follows"`
	Experiences     []*Experience `json:"Experiences" gorm:"many2many:user_experiences"`
	Educations      []*Education  `json:"Educations" gorm:"many2many:user_educations"`
}

type Reset struct {
	ID     string `json:"id" gorm:"type:varchar(191)"`
	UserId string `gorm:"reference:User;type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User   *User  `json:"User"`
}

type Message struct {
	ID        string    `json:"id" gorm:"type:varchar(191)"`
	Text      string    `json:"Text"`
	User1Id   string    `gorm:"reference:User;type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User1     *User     `json:"User1"`
	User2Id   string    `gorm:"reference:User;type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User2     *User     `json:"User2"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type UserFollow struct {
	ID       string `json:"id" gorm:"type:varchar(191)"`
	UserId   string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FollowId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserVisit struct {
	ID      string `json:"id" gorm:"type:varchar(191)"`
	UserId  string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	VisitId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserExperience struct {
	ID           string `json:"id" gorm:"type:varchar(191)"`
	UserId       string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ExperienceId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserEducation struct {
	ID          string `json:"id" gorm:"type:varchar(191)"`
	UserId      string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EducationId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Education struct {
	ID        string `json:"id" gorm:"type:varchar(191)"`
	School    string `json:"School"`
	Field     string `json:"Field"`
	StartedAt string `json:"StartedAt"`
	EndedAt   string `json:"EndedAt"`
}

type Experience struct {
	ID        string `json:"id" gorm:"type:varchar(191)"`
	Position  string `json:"Position"`
	Desc      string `json:"Desc"`
	Company   string `json:"Company"`
	StartedAt string `json:"StartedAt"`
	EndedAt   string `json:"EndedAt"`
	IsActive  bool   `json:"IsActive"`
}

type Activation struct {
	ID     string `json:"ID"`
	UserId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User   *User  `json:"User"`
}

type Connection struct {
	ID      string `json:"id" gorm:"type:varchar(191)"`
	User1ID string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User1   *User  `json:"User1"`
	User2ID string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User2   *User  `json:"User2"`
}

type ConnectRequest struct {
	ID      string `json:"id" gorm:"type:varchar(191)"`
	User1ID string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User1   *User  `json:"User1"`
	User2ID string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User2   *User  `json:"User2"`
}

type InputLogin struct {
	Email    string `json:"email"`
	Password string `json:"passowrd"`
}

type InputRegister struct {
	Email     string `json:"Email"`
	Password  string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	MidName   string `json:"MidName"`
}

type ConnectStatus string

const (
	ConnectStatusConnected    ConnectStatus = "Connected"
	ConnectStatusSentByUser1  ConnectStatus = "SentByUser1"
	ConnectStatusSentByUser2  ConnectStatus = "SentByUser2"
	ConnectStatusNotConnected ConnectStatus = "NotConnected"
)

var AllConnectStatus = []ConnectStatus{
	ConnectStatusConnected,
	ConnectStatusSentByUser1,
	ConnectStatusSentByUser2,
	ConnectStatusNotConnected,
}

func (e ConnectStatus) IsValid() bool {
	switch e {
	case ConnectStatusConnected, ConnectStatusSentByUser1, ConnectStatusSentByUser2, ConnectStatusNotConnected:
		return true
	}
	return false
}

func (e ConnectStatus) String() string {
	return string(e)
}

func (e *ConnectStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ConnectStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ConnectStatus", str)
	}
	return nil
}

func (e ConnectStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type InputEducation struct {
	School    string `json:"School"`
	Field     string `json:"Field"`
	StartedAt string `json:"StartedAt"`
	EndedAt   string `json:"EndedAt"`
}

type InputExperience struct {
	Position  string `json:"Position"`
	Desc      string `json:"Desc"`
	Company   string `json:"Company"`
	StartedAt string `json:"StartedAt"`
	EndedAt   string `json:"EndedAt"`
	IsActive  bool   `json:"IsActive"`
}

type InputMessage struct {
	Text    string `json:"Text"`
	User1Id string `json:"User1Id"`
	User2Id string `json:"User2Id"`
}

type InputUser struct {
	FirstName       string `json:"FirstName"`
	LastName        string `json:"LastName"`
	MidName         string `json:"MidName"`
	ProfilePhoto    string `json:"ProfilePhoto"`
	BackgroundPhoto string `json:"BackgroundPhoto"`
	Headline        string `json:"Headline"`
	Pronoun         string `json:"Pronoun"`
	About           string `json:"About"`
	Location        string `json:"Location"`
}
