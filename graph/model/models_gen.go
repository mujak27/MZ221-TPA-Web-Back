// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type TypeConnection struct {
	ConnectionStatus ConnectStatus `json:"connectionStatus"`
	Text             string        `json:"text"`
}

type EnumMessageType string

const (
	EnumMessageTypeText      EnumMessageType = "text"
	EnumMessageTypeVideoCall EnumMessageType = "videoCall"
	EnumMessageTypePost      EnumMessageType = "post"
	EnumMessageTypeUser      EnumMessageType = "user"
)

var AllEnumMessageType = []EnumMessageType{
	EnumMessageTypeText,
	EnumMessageTypeVideoCall,
	EnumMessageTypePost,
	EnumMessageTypeUser,
}

func (e EnumMessageType) IsValid() bool {
	switch e {
	case EnumMessageTypeText, EnumMessageTypeVideoCall, EnumMessageTypePost, EnumMessageTypeUser:
		return true
	}
	return false
}

func (e EnumMessageType) String() string {
	return string(e)
}

func (e *EnumMessageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EnumMessageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid enumMessageType", str)
	}
	return nil
}

func (e EnumMessageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
