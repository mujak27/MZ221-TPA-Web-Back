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

// Post is the resolver for the Post field.
func (r *commentResolver) Post(ctx context.Context, obj *model.Comment) (*model.Post, error) {
	var post *model.Post
	if err := r.DB.First(&post, "id = ?", obj.PostId).Error; err != nil {
		return nil, err
	}
	return post, nil
}

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

// Likes is the resolver for the Likes field.
func (r *commentResolver) Likes(ctx context.Context, obj *model.Comment) ([]*model.User, error) {
	var commentLikes []*model.CommentLike
	if err := r.DB.Find(&commentLikes, "comment_id = ?", obj.ID).Error; err != nil {
		return nil, err
	}

	userIds := lo.Map[*model.CommentLike, string](commentLikes, func(x *model.CommentLike, _ int) string {
		return x.LikeId
	})

	return UsersById(r.Resolver, userIds)
}

// CreatedAt is the resolver for the CreatedAt field.
func (r *commentResolver) CreatedAt(ctx context.Context, obj *model.Comment) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// CreatePost is the resolver for the CreatePost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.InputPost) (*model.Post, error) {
	userId := auth.JwtGetValue(ctx).Userid
	fmt.Println(userId)

	user, err := UserById(r.Resolver, userId)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		ID:             uuid.NewString(),
		Text:           input.Text,
		SenderId:       userId,
		Comments:       []*model.Comment{},
		AttachmentLink: input.AttachmentLink,
	}
	r.posts = append(r.posts, post)

	r.DB.Create(post)

	activityText := user.FirstName + " " + user.LastName + " has created a new post " + post.Text
	AddActivity(r.Resolver, userId, activityText)

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
	like, _ := r.Query().IsLikePost(ctx, id)
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

// LikeComment is the resolver for the LikeComment field.
func (r *mutationResolver) LikeComment(ctx context.Context, id string) (interface{}, error) {
	myId := getId(ctx)
	like, _ := r.Query().IsLikeComment(ctx, id)
	if like {
		return map[string]interface{}{
			"status": "have been liked",
		}, nil
	}

	commentLike := &model.CommentLike{
		ID:        uuid.NewString(),
		CommentId: id,
		LikeId:    myId,
	}
	r.comment_likes = append(r.comment_likes, commentLike)
	r.DB.Create(commentLike)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// UnLikeComment is the resolver for the UnLikeComment field.
func (r *mutationResolver) UnLikeComment(ctx context.Context, id string) (interface{}, error) {
	myId := getId(ctx)
	var commentLike *model.CommentLike
	err := r.DB.First(&commentLike, "like_id = ? and comment_id = ?", myId, id).Error
	if err != nil {
		return map[string]interface{}{
			"status": "not liked",
		}, err
	}
	r.DB.Delete(commentLike)
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

// CreatedAt is the resolver for the CreatedAt field.
func (r *postResolver) CreatedAt(ctx context.Context, obj *model.Post) (string, error) {
	panic(fmt.Errorf("not implemented"))
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

	idList, err := getFriendIds(r.Resolver, ctx)
	if err != nil {
		return nil, err
	}

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
	idList, err := getFriendIds(r.Resolver, ctx)
	if err != nil {
		return nil, err
	}

	if keyword == "" {
		keyword = "%"
	}

	var posts []*model.Post
	if err := r.DB.Limit(limit).Offset(offset).Find(&posts, "sender_id IN ? and text like ?", idList, "%"+keyword+"%").Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// Comment is the resolver for the Comment field.
func (r *queryResolver) Comment(ctx context.Context, id string) (*model.Comment, error) {
	var comment *model.Comment
	var err error
	if err = r.DB.Find(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// Comments is the resolver for the Comments field.
func (r *queryResolver) Comments(ctx context.Context, commentID *string, postID string) ([]*model.Comment, error) {
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

// IsLikePost is the resolver for the IsLikePost field.
func (r *queryResolver) IsLikePost(ctx context.Context, id string) (bool, error) {
	myId := getId(ctx)

	var postLike *model.PostLike
	err := r.DB.First(&postLike, "like_id = ? and post_id = ?", myId, id).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

// IsLikeComment is the resolver for the IsLikeComment field.
func (r *queryResolver) IsLikeComment(ctx context.Context, id string) (bool, error) {
	myId := getId(ctx)

	var commentLike *model.CommentLike
	err := r.DB.First(&commentLike, "like_id = ? and comment_id = ?", myId, id).Error
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

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *postResolver) AttachmentLink(ctx context.Context, obj *model.Post) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
