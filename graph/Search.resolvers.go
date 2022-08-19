package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"
)

// Search is the resolver for the Search field.
func (r *queryResolver) Search(ctx context.Context, keyword string) (*model.Search, error) {
	panic(fmt.Errorf("not implemented"))
}

// Users is the resolver for the Users field.
func (r *searchResolver) Users(ctx context.Context, obj *model.Search) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Post is the resolver for the Post field.
func (r *searchResolver) Post(ctx context.Context, obj *model.Search) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// Search returns generated.SearchResolver implementation.
func (r *Resolver) Search() generated.SearchResolver { return &searchResolver{r} }

type searchResolver struct{ *Resolver }
