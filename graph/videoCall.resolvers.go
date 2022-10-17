package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"
	"time"
)

// OfferCandidates is the resolver for the OfferCandidates field.
func (r *mutationResolver) OfferCandidates(ctx context.Context, id string, input model.InputCandidates) (interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}

// AnswerCandidates is the resolver for the AnswerCandidates field.
func (r *mutationResolver) AnswerCandidates(ctx context.Context, id string, input model.InputCandidates) (interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}

// HangUp is the resolver for the HangUp field.
func (r *mutationResolver) HangUp(ctx context.Context, id string) (interface{}, error) {
	panic(fmt.Errorf("not implemented"))
}

// MutationTestSubs is the resolver for the mutationTestSubs field.
func (r *mutationResolver) MutationTestSubs(ctx context.Context, str string) (string, error) {
	vidcall := model.VideoCall{
		ID:               str,
		User1:            &model.User{},
		User2:            &model.User{},
		OfferCandidates:  []string{},
		AnswerCandidates: []string{},
	}
	for idx, ch := range channels {
		select {
		case ch <- &vidcall:
		default:
			fmt.Printf("%d closed\n", idx)
			// return "closed", nil
		}
	}
	return str, nil
}

// VideoCallStatus is the resolver for the VideoCallStatus field.
func (r *queryResolver) VideoCallStatus(ctx context.Context, id string) (model.VideoCallStatus, error) {
	panic(fmt.Errorf("not implemented"))
}

// SubscribeTest is the resolver for the subscribeTest field.
func (r *subscriptionResolver) SubscribeTest(ctx context.Context) (<-chan *model.VideoCall, error) {
	ch := make(chan *model.VideoCall)
	channels = append(channels, ch)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			vidcall := model.VideoCall{
				ID:               time.Now().String(),
				User1:            &model.User{},
				User2:            &model.User{},
				OfferCandidates:  []string{},
				AnswerCandidates: []string{},
			}
			select {
			case ch <- &vidcall:
			default:
				fmt.Println(("closed"))
				return
			}
		}
	}()
	return ch, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var channels = []chan *model.VideoCall{}
