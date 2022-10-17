package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// User1 is the resolver for the User1 field.
func (r *messageResolver) User1(ctx context.Context, obj *model.Message) (*model.User, error) {
	return UserById(r.Resolver, obj.User1Id)
}

// User2 is the resolver for the User2 field.
func (r *messageResolver) User2(ctx context.Context, obj *model.Message) (*model.User, error) {
	return UserById(r.Resolver, obj.User2Id)
}

// CreatedAt is the resolver for the CreatedAt field.
func (r *messageResolver) CreatedAt(ctx context.Context, obj *model.Message) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// SendMessage is the resolver for the SendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.InputMessage) (model.MutationStatus, error) {
	message := &model.Message{
		ID:          uuid.NewString(),
		Text:        input.Text,
		ImageLink:   input.ImageLink,
		MessageType: input.MessageType,
		User1Id:     input.User1Id,
		User2Id:     input.User2Id,
		CreatedAt:   time.Time{},
	}
	r.messages = append(r.messages, message)
	r.DB.Create(message)

	return model.MutationStatusSuccess, nil
}

// Messages is the resolver for the Messages field.
func (r *queryResolver) Messages(ctx context.Context, id1 string, id2 string) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.DB.Find(&messages, "(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?) order by created_at", id1, id2, id2, id1).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// RecentMessage is the resolver for the RecentMessage field.
func (r *queryResolver) RecentMessage(ctx context.Context) ([]*model.Message, error) {
	myId := getId(ctx)
	var messages []*model.Message
	if err := r.DB.Where("id in (	select 		id	from messages as m2	where 		(user1_id = ? or user2_id = ?)		and created_at = (			select 				max(created_at)			from messages as m3			where 				(user1_id = m2.user1_id and user2_id = m2.user2_id)				or (user1_id = m2.user2_id and user2_id = m2.user1_id)		))", myId, myId).Order("created_at desc").Find(&messages).Error; err != nil {
		fmt.Println("error")
		return nil, err
	}

	blockIds, _ := getBlockIds(r.Resolver, ctx)

	messages = lo.Filter(messages, func(x *model.Message, _ int) bool {
		return !lo.Contains(blockIds, x.User1Id) && !lo.Contains(blockIds, x.User2Id)
	})

	fmt.Println(messages)

	return messages, nil
}

// GetMessages is the resolver for the getMessages field.
func (r *subscriptionResolver) GetMessages(ctx context.Context, id string) (<-chan []*model.Message, error) {
	ch := make(chan []*model.Message)

	id1 := getId(ctx)
	id2 := id

	var messages []*model.Message
	if err := r.DB.Find(&messages, "(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?) order by created_at", id1, id2, id2, id1).Error; err != nil {
		return nil, err
	}
	ch <- messages
	return ch, nil
}

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type messageResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
