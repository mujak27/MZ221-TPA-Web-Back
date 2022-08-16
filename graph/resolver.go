package graph

import (
	"MZ221-TPA-Web-Back/graph/model"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//go:generate go run github.com/99designs/gqlgen generate

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	posts []*model.Post
	users []*model.User
	DB    *gorm.DB
}
