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

// User is the resolver for the User field.
func (r *activationResolver) User(ctx context.Context, obj *model.Activation) (*model.User, error) {
	return UserById(r.Resolver, obj.UserId)
}

// User1 is the resolver for the User1 field.
func (r *connectRequestResolver) User1(ctx context.Context, obj *model.ConnectRequest) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// User2 is the resolver for the User2 field.
func (r *connectRequestResolver) User2(ctx context.Context, obj *model.ConnectRequest) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// User1 is the resolver for the User1 field.
func (r *connectionResolver) User1(ctx context.Context, obj *model.Connection) (*model.User, error) {
	return UserById(r.Resolver, obj.User1ID)
}

// User2 is the resolver for the User2 field.
func (r *connectionResolver) User2(ctx context.Context, obj *model.Connection) (*model.User, error) {
	return UserById(r.Resolver, obj.User2ID)
}

// User1 is the resolver for the User1 field.
func (r *followResolver) User1(ctx context.Context, obj *model.Follow) (*model.User, error) {
	return UserById(r.Resolver, obj.User1ID)
}

// User2 is the resolver for the User2 field.
func (r *followResolver) User2(ctx context.Context, obj *model.Follow) (*model.User, error) {
	return UserById(r.Resolver, obj.User2ID)
}

// SendActivation is the resolver for the SendActivation field.
func (r *mutationResolver) SendActivation(ctx context.Context, id string) (interface{}, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	SendActivationLink(r.Resolver, user)
	return user, nil
}

// Activate is the resolver for the Activate field.
func (r *mutationResolver) Activate(ctx context.Context, id string) (interface{}, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user.IsActive = true
	r.DB.Save(user)
	return map[string]interface{}{
		"status": "success",
	}, nil
}

// Follow is the resolver for the Follow field.
func (r *mutationResolver) Follow(ctx context.Context, id1 string, id2 string) (interface{}, error) {
	var follow *model.Follow
	err := r.DB.First(&follow, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err == nil {
		return map[string]interface{}{
			"status": "already exist",
		}, nil
	}

	follow = &model.Follow{
		ID:      uuid.NewString(),
		User1ID: id1,
		User2ID: id2,
	}
	r.follows = append(r.follows, follow)

	r.DB.Create(follow)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// UnFollow is the resolver for the UnFollow field.
func (r *mutationResolver) UnFollow(ctx context.Context, id1 string, id2 string) (interface{}, error) {
	var follow *model.Follow
	err := r.DB.First(&follow, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err != nil {
		return map[string]interface{}{
			"status": "failed",
		}, nil
	}
	r.DB.Delete(follow)
	return map[string]interface{}{
		"status": "success",
	}, nil
}

// SendConnectRequest is the resolver for the SendConnectRequest field.
func (r *mutationResolver) SendConnectRequest(ctx context.Context, id1 string, id2 string) (interface{}, error) {
	var connectRequest *model.ConnectRequest
	err := r.DB.First(&connectRequest, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err == nil {
		return map[string]interface{}{
			"status": "already exists",
		}, nil
	}
	connectRequest = &model.ConnectRequest{
		ID:      uuid.NewString(),
		User1ID: id1,
		User2ID: id2,
	}

	r.connectRequests = append(r.connectRequests, connectRequest)

	r.DB.Create(connectRequest)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// DeleteConnectRequest is the resolver for the DeleteConnectRequest field.
func (r *mutationResolver) DeleteConnectRequest(ctx context.Context, id1 string, id2 string) (interface{}, error) {
	var connectRequest *[]model.ConnectRequest
	err1 := r.DB.Find(&connectRequest, "user1_id = ? and user2_id = ?", id1, id2).Error
	r.DB.Delete(connectRequest)
	err2 := r.DB.Find(&connectRequest, "user1_id = ? and user2_id = ?", id2, id1).Error
	r.DB.Delete(connectRequest)
	if err1 != nil && err2 != nil {
		return map[string]interface{}{
			"status": "not found",
		}, err1
	}
	// r.DB.Delete(connectRequest)
	return map[string]interface{}{
		"status": "success",
	}, nil
}

// AcceptConnectRequest is the resolver for the AcceptConnectRequest field.
func (r *mutationResolver) AcceptConnectRequest(ctx context.Context, id1 string, id2 string) (interface{}, error) {
	status, err := r.DeleteConnectRequest(ctx, id1, id2)
	if err != nil {
		return status, err
	}

	id1, id2 = SortIdAsc(id1, id2)

	connection := &model.Connection{
		ID:      uuid.NewString(),
		User1ID: id1,
		User2ID: id2,
	}
	r.connections = append(r.connections, connection)
	r.DB.Create(connection)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// UnConnect is the resolver for the UnConnect field.
func (r *mutationResolver) UnConnect(ctx context.Context, id1 string, id2 string) (interface{}, error) {

	id1, id2 = SortIdAsc(id1, id2)

	var connection *model.Connection
	err := r.DB.First(&connection, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err != nil {
		return map[string]interface{}{
			"status": "not found",
		}, err
	}
	r.DB.Delete(connection)
	return map[string]interface{}{
		"status": "success",
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return UserById(r.Resolver, id)
}

// UsersByName is the resolver for the UsersByName field.
func (r *queryResolver) UsersByName(ctx context.Context, name *string) ([]*model.User, error) {
	var users []*model.User
	if err := r.DB.Find(&users, "concat(first_name, mid_name, last_name) like ?", "%"+*name+"%").Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Login is the resolver for the login field.
func (r *queryResolver) Login(ctx context.Context, input model.InputLogin) (interface{}, error) {
	var user *model.User
	if err := r.DB.Where("email = ?", input.Email).First(&user, "password = ?", input.Password).Error; err != nil {
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

	r.DB.Create(newUser)

	SendActivationLink(r.Resolver, newUser)

	token, err := auth.JwtGenerate(ctx, newUser.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

// Activation is the resolver for the Activation field.
func (r *queryResolver) Activation(ctx context.Context, id string) (*model.Activation, error) {
	var activation *model.Activation
	if err := r.DB.First(&activation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return activation, nil
}

// IsFollow is the resolver for the IsFollow field.
func (r *queryResolver) IsFollow(ctx context.Context, id1 string, id2 string) (bool, error) {
	var follow *model.Follow
	if err := r.DB.First(&follow, "user1_id = ? and user2_id = ?", id1, id2).Error; err != nil {
		return false, nil
	}
	return true, nil
}

// IsConnect is the resolver for the IsConnect field.
func (r *queryResolver) IsConnect(ctx context.Context, id1 string, id2 string) (model.ConnectStatus, error) {

	sortedId1, sortedId2 := SortIdAsc(id1, id2)

	var connection *model.Connection
	err := r.DB.First(&connection, "user1_id = ? and user2_id = ?", sortedId1, sortedId2).Error
	if err == nil {
		return model.ConnectStatusConnected, nil
	}

	var connectRequest *model.ConnectRequest
	err = r.DB.First(&connectRequest, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err == nil {
		return model.ConnectStatusSentByUser1, nil
	}

	err = r.DB.First(&connectRequest, "user1_id = ? and user2_id = ?", id2, id1).Error
	if err == nil {
		return model.ConnectStatusSentByUser2, nil
	}

	return model.ConnectStatusNotConnected, nil

}

// User1 is the resolver for the User1 field.
func (r *visitResolver) User1(ctx context.Context, obj *model.Visit) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// User2 is the resolver for the User2 field.
func (r *visitResolver) User2(ctx context.Context, obj *model.Visit) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Activation returns generated.ActivationResolver implementation.
func (r *Resolver) Activation() generated.ActivationResolver { return &activationResolver{r} }

// ConnectRequest returns generated.ConnectRequestResolver implementation.
func (r *Resolver) ConnectRequest() generated.ConnectRequestResolver {
	return &connectRequestResolver{r}
}

// Connection returns generated.ConnectionResolver implementation.
func (r *Resolver) Connection() generated.ConnectionResolver { return &connectionResolver{r} }

// Follow returns generated.FollowResolver implementation.
func (r *Resolver) Follow() generated.FollowResolver { return &followResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Visit returns generated.VisitResolver implementation.
func (r *Resolver) Visit() generated.VisitResolver { return &visitResolver{r} }

type activationResolver struct{ *Resolver }
type connectRequestResolver struct{ *Resolver }
type connectionResolver struct{ *Resolver }
type followResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type visitResolver struct{ *Resolver }
