package graph

import (
	"MZ221-TPA-Web-Back/graph/model"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen generate

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	experiences      []*model.Experience
	educations       []*model.Education
	connectRequests  []*model.ConnectRequest
	connections      []*model.Connection
	activations      []*model.Activation
	posts            []*model.Post
	postLikes        []*model.PostLike
	comments         []*model.Comment
	comment_likes    []*model.CommentLike
	users            []*model.User
	user_follows     []*model.UserFollow
	user_visits      []*model.UserVisit
	user_experiences []*model.UserExperience
	user_educations  []*model.UserEducation
	activities       []*model.Activity
	messages         []*model.Message
	resets           []*model.Reset
	DB               *gorm.DB
	// follows         []*model.Follow
}
