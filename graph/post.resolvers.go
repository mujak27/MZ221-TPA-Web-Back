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

// Sender is the resolver for the Sender field.
func (r *commentResolver) Sender(ctx context.Context, obj *model.Comment) (*model.User, error) {
	return UserById(r.Resolver, obj.SenderId)
}

// Replies is the resolver for the Replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *model.Comment) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := r.DB.Find(&comments, "replied_to_id = ?", obj.ID).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// CreatePost is the resolver for the CreatePost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.InputPost) (*model.Post, error) {
	userId := auth.JwtGetValue(ctx).Userid
	fmt.Println(userId)

	post := &model.Post{
		ID:       uuid.NewString(),
		Text:     input.Text,
		SenderId: userId,
		Comments: []*model.Comment{},
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

// LikePost is the resolver for the LikePost field.
func (r *mutationResolver) LikePost(ctx context.Context, id string) (interface{}, error) {
	myId := getId(ctx)
	like, _ := r.Query().IsLiked(ctx, id)
	if like {
		return map[string]interface{}{
			"status": "have been liked",
		}, nil
	}

	postLike := &model.PostLike{
		ID:     uuid.NewString(),
		PostId: id,
		LikeId: myId,
	}
	r.postLikes = append(r.postLikes, postLike)
	r.DB.Create(postLike)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// UnLikePost is the resolver for the UnLikePost field.
func (r *mutationResolver) UnLikePost(ctx context.Context, id string) (interface{}, error) {
	myId := getId(ctx)
	var postLike *model.PostLike
	err := r.DB.First(&postLike, "like_id = ? and post_id = ?", myId, id).Error
	if err != nil {
		return map[string]interface{}{
			"status": "not liked",
		}, err
	}
	r.DB.Delete(postLike)
	return map[string]interface{}{
		"status": "success",
	}, nil
}

// CommentPost is the resolver for the CommentPost field.
func (r *mutationResolver) CommentPost(ctx context.Context, input *model.InputComment) (interface{}, error) {
	myId := getId(ctx)
	comment := &model.Comment{
		ID:       uuid.NewString(),
		Text:     input.Text,
		SenderId: myId,
		PostId:   input.PostID,
	}
	if input.RepliedToID != "" {
		comment.RepliedToId = &input.RepliedToID
	} else {
	}
	r.comments = append(r.comments, comment)
	if err := r.DB.Create(comment).Error; err != nil {
		return map[string]interface{}{
			"status": "failed",
		}, err
	}

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// Sender is the resolver for the Sender field.
func (r *postResolver) Sender(ctx context.Context, obj *model.Post) (*model.User, error) {
	return UserById(r.Resolver, obj.SenderId)
}

// Comments is the resolver for the Comments field.
func (r *postResolver) Comments(ctx context.Context, obj *model.Post) ([]*model.Comment, error) {
	var comments []*model.Comment
	if err := r.DB.Find(&comments, "post_id = ? and replied_to_id is NULL", obj.ID).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// Likes is the resolver for the Likes field.
func (r *postResolver) Likes(ctx context.Context, obj *model.Post) ([]*model.User, error) {
	var postLikes []*model.PostLike
	if err := r.DB.Find(&postLikes, "post_id = ?", obj.ID).Error; err != nil {
		return nil, err
	}

	userIds := lo.Map[*model.PostLike, string](postLikes, func(x *model.PostLike, _ int) string {
		return x.LikeId
	})
	fmt.Println(userIds)

	return UsersById(r.Resolver, userIds)
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

	if err := r.DB.Find(&follows, "user_id = ?", myId).Error; err != nil {
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

	fmt.Println("----------------")
	fmt.Println(idList)

	var posts []*model.Post
	if err := r.DB.Limit(limit).Offset(offset).Find(&posts, "sender_id IN ?", idList).Error; err != nil {
		return nil, err
	}
	fmt.Println(posts)

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

// CommentReplies is the resolver for the CommentReplies field.
func (r *queryResolver) CommentReplies(ctx context.Context, commentID *string, postID string) ([]*model.Comment, error) {
	var comments []*model.Comment
	var err error
	if commentID != nil {
		if err = r.DB.Find(&comments, "replied_to_id = ?", commentID).Error; err != nil {
			return nil, err
		}
	} else {
		if err = r.DB.Find(&comments, "replied_to_id is NULL and post_id = ?", postID).Error; err != nil {
			return nil, err
		}
	}
	return comments, nil
}

// IsLiked is the resolver for the IsLiked field.
func (r *queryResolver) IsLiked(ctx context.Context, id string) (bool, error) {
	myId := getId(ctx)

	var postLike *model.PostLike
	err := r.DB.First(&postLike, "like_id = ? and post_id = ?", myId, id).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
