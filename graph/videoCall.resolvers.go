package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"
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

// VideoCallStatus is the resolver for the VideoCallStatus field.
func (r *queryResolver) VideoCallStatus(ctx context.Context, id string) (model.VideoCallStatus, error) {
	panic(fmt.Errorf("not implemented"))
}
