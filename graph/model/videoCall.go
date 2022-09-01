package model

import (
	"fmt"
	"io"
	"strconv"
)

type InputCandidates struct {
	UserID string `json:"userId"`
}

type VideoCall struct {
	ID               string   `json:"ID"`
	User1            *User    `json:"User1"`
	User2            *User    `json:"User2"`
	OfferCandidates  []string `json:"OfferCandidates"`
	AnswerCandidates []string `json:"AnswerCandidates"`
}

type VideoCallStatus string

const (
	VideoCallStatusNull         VideoCallStatus = "Null"
	VideoCallStatusWaitingUser1 VideoCallStatus = "WaitingUser1"
	VideoCallStatusWaitingUser2 VideoCallStatus = "WaitingUser2"
	VideoCallStatusOnGoing      VideoCallStatus = "OnGoing"
)

var AllVideoCallStatus = []VideoCallStatus{
	VideoCallStatusNull,
	VideoCallStatusWaitingUser1,
	VideoCallStatusWaitingUser2,
	VideoCallStatusOnGoing,
}

func (e VideoCallStatus) IsValid() bool {
	switch e {
	case VideoCallStatusNull, VideoCallStatusWaitingUser1, VideoCallStatusWaitingUser2, VideoCallStatusOnGoing:
		return true
	}
	return false
}

func (e VideoCallStatus) String() string {
	return string(e)
}

func (e *VideoCallStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = VideoCallStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid VideoCallStatus", str)
	}
	return nil
}

func (e VideoCallStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
