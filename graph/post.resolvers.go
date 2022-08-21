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
	fmt.Println(userId)

	post := &model.Post{
		ID:       uuid.NewString(),
		Text:     input.Text,
		SenderId: userId,
	}
	r.posts = append(r.posts, post)

	r.DB.Create(post)

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

// Sender is the resolver for the Sender field.
func (r *postResolver) Sender(ctx context.Context, obj *model.Post) (*model.User, error) {
	return UserById(r.Resolver, obj.SenderId)
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
func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]*model.Post, error) {
	fmt.Println(limit, offset)

	var idList []string
	myId := auth.JwtGetValue(ctx).Userid
	idList = append(idList, myId)

	var follows []*model.UserFollow

	if err := r.DB.Find(&follows, "user1_id = ?", myId).Error; err != nil {
		return nil, err
	}

	followIds := lo.Map[*model.UserFollow, string](follows, func(x *model.UserFollow, _ int) string {
		return x.FollowId
	})
	idList = append(idList, followIds...)

	var connections []*model.Connection

	if err := r.DB.Find(&connections, "user1_id = ?", myId).Error; err != nil {
		return nil, err
	}
	connectionIds := lo.Map[*model.Connection, string](connections, func(x *model.Connection, _ int) string {
		return x.User2ID
	})
	idList = append(idList, connectionIds...)

	if err := r.DB.Find(&connections, "user2_id = ?", myId).Error; err != nil {
		return nil, err
	}
	connectionIds = lo.Map[*model.Connection, string](connections, func(x *model.Connection, _ int) string {
		return x.User1ID
	})
	idList = append(idList, connectionIds...)

	idList = lo.Uniq[string](idList)

	fmt.Println(idList)

	var posts []*model.Post
	if err := r.DB.Limit(limit).Offset(offset).Find(&posts, "sender_id IN ?", idList).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// PostsByUserID is the resolver for the PostsByUserId field.
func (r *queryResolver) PostsByUserID(ctx context.Context, id string) ([]*model.Post, error) {
	panic(fmt.Errorf("not implemented"))
}

// PostsByKeyword is the resolver for the PostsByKeyword field.
func (r *queryResolver) PostsByKeyword(ctx context.Context, keyword string, limit int, offset int) ([]*model.Post, error) {
	if keyword == "" {
		keyword = "%"
	}

	var idList []string
	myId := auth.JwtGetValue(ctx).Userid
	idList = append(idList, myId)

	user, err := UserById(r.Resolver, myId)
	if err != nil {
		return nil, err
	}

	followIds := lo.Map[*model.User, string](user.Follows, func(x *model.User, _ int) string {
		return x.ID
	})
	idList = append(idList, followIds...)

	var connections []*model.Connection
	if err := r.DB.Find(&connections, "user1_id = ?", myId).Error; err != nil {
		return nil, err
	}
	connectionIds := lo.Map[*model.Connection, string](connections, func(x *model.Connection, _ int) string {
		return x.User2ID
	})
	idList = append(idList, connectionIds...)

	if err := r.DB.Find(&connections, "user2_id = ?", myId).Error; err != nil {
		return nil, err
	}
	connectionIds = lo.Map[*model.Connection, string](connections, func(x *model.Connection, _ int) string {
		return x.User1ID
	})
	idList = append(idList, connectionIds...)

	idList = lo.Uniq[string](idList)

	fmt.Println(idList)

	var posts []*model.Post
	if err := r.DB.Limit(limit).Offset(offset).Find(&posts, "sender_id IN ? and text like ?", idList, "%"+keyword+"%").Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
