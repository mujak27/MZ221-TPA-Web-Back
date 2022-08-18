package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/auth"
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// CreatePost is the resolver for the CreatePost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.InputPost) (*model.Post, error) {
	userId := auth.JwtGetValue(ctx).Userid

	post := &model.Post{
		ID:   uuid.NewString(),
		Text: input.Text,
	}
	r.posts = append(r.posts, post)

	sender := &model.PostSender{
		ID:     uuid.NewString(),
		UserId: userId,
		PostId: post.ID,
	}

	r.senders = append(r.senders, sender)

	r.DB.Create(post)
	r.DB.Create(sender)

	return post, nil
}

// UpdatePost is the resolver for the UpdatePost field.
func (r *mutationResolver) UpdatePost(ctx context.Context, id string) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// DeletePost is the resolver for the DeletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// AddPostSender is the resolver for the AddPostSender field.
func (r *mutationResolver) AddPostSender(ctx context.Context, id string) (*model.Post, error) {
	var post *model.Post
	if err := r.DB.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}

	userId := auth.JwtGetValue(ctx).Userid

	sender := &model.PostSender{
		ID:     uuid.NewString(),
		UserId: userId,
		PostId: id,
	}
	r.senders = append(r.senders, sender)
	r.DB.Create(sender)

	return post, nil
}

// Senders is the resolver for the Senders field.
func (r *postResolver) Senders(ctx context.Context, obj *model.Post) ([]*model.User, error) {
	var postSenders []*model.PostSender
	if err := r.DB.Find(&postSenders, obj.Senders).Error; err != nil {
		return nil, err
	}

	fmt.Println(postSenders)
	var sendersId = lo.Map[*model.PostSender, string](postSenders, func(x *model.PostSender, _ int) string {
		return x.UserId
	})

	fmt.Println(sendersId)

	var users []*model.User
	if err := r.DB.Find(&users, sendersId).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Post is the resolver for the Post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	var post *model.Post
	if err := r.DB.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return post, nil
}

// Posts is the resolver for the Posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// PostsByUserID is the resolver for the PostsByUserId field.
func (r *queryResolver) PostsByUserID(ctx context.Context, id string) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
