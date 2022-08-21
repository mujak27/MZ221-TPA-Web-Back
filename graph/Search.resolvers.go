package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"

	"github.com/samber/lo"
)

// Search is the resolver for the Search field.
func (r *queryResolver) Search(ctx context.Context, keyword string, limit int, offset int) (*model.Search, error) {
	var search *model.Search
	search = &model.Search{}

	users, err := r.UsersByName(ctx, &keyword, limit, offset)
	if err != nil {
		return nil, err
	}
	search.Users = users

	posts, _ := r.PostsByKeyword(ctx, keyword, limit, offset)
	search.Posts = posts

	return search, nil
}

// Users is the resolver for the Users field.
func (r *searchResolver) Users(ctx context.Context, obj *model.Search) ([]*model.User, error) {
	userIds := lo.Map[*model.User, string](obj.Users, func(x *model.User, _ int) string {
		return x.ID
	})

	var users []*model.User
	if err := r.DB.Find(&users, userIds).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Posts is the resolver for the Posts field.
func (r *searchResolver) Posts(ctx context.Context, obj *model.Search) ([]*model.Post, error) {
	postIds := lo.Map[*model.Post, string](obj.Posts, func(x *model.Post, _ int) string {
		return x.ID
	})

	var posts []*model.Post
	if err := r.DB.Find(&posts, postIds).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// Search returns generated.SearchResolver implementation.
func (r *Resolver) Search() generated.SearchResolver { return &searchResolver{r} }

type searchResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *searchResolver) Post(ctx context.Context, obj *model.Search) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}
