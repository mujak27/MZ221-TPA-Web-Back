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
)

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Login is the resolver for the login field.
func (r *queryResolver) Login(ctx context.Context, input model.InputLogin) (interface{}, error) {
	var user *model.User
	if err := r.DB.Where("email = ?", input.Email).First(&user, "password = ?", input.Passowrd).Error; err != nil {
		return nil, err
	}
	fmt.Println(user)

	token, err := auth.JwtGenerate(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

// Register is the resolver for the register field.
func (r *queryResolver) Register(ctx context.Context, input *model.InputRegister) (interface{}, error) {
	newUser := &model.User{
		ID:              uuid.NewString(),
		Email:           input.Email,
		Password:        input.Password,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		MidName:         input.MidName,
		IsActive:        false,
		ProfilePhoto:    "",
		BackgroundPhoto: "",
		Headline:        "",
		Pronoun:         "",
		ProfileLink:     "",
		About:           "",
		Location:        "",
	}
	r.users = append(r.users, newUser)

	token, err := auth.JwtGenerate(ctx, newUser.ID)
	if err != nil {
		return nil, err
	}
	r.DB.Create(newUser)

	return map[string]interface{}{
		"token": token,
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
