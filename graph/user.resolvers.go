package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/auth"
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// User is the resolver for the User field.
func (r *activationResolver) User(ctx context.Context, obj *model.Activation) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
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
	panic(fmt.Errorf("not implemented"))
}

// User2 is the resolver for the User2 field.
func (r *connectionResolver) User2(ctx context.Context, obj *model.Connection) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// User1 is the resolver for the User1 field.
func (r *messageResolver) User1(ctx context.Context, obj *model.Message) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// User2 is the resolver for the User2 field.
func (r *messageResolver) User2(ctx context.Context, obj *model.Message) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// CreatedAt is the resolver for the CreatedAt field.
func (r *messageResolver) CreatedAt(ctx context.Context, obj *model.Message) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// UpdateProfile is the resolver for the UpdateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, input model.InputUser) (*model.User, error) {
	myId := getId(ctx)

	user, err := UserById(r.Resolver, myId)
	if err != nil {
		return nil, err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.MidName = input.MidName
	user.ProfilePhoto = input.ProfilePhoto
	user.BackgroundPhoto = input.BackgroundPhoto
	user.Headline = input.Headline
	user.Pronoun = input.Pronoun
	user.About = input.About
	user.Location = input.Location

	r.DB.Save(user)

	return user, nil
}

// ForgetPassword is the resolver for the ForgetPassword field.
func (r *mutationResolver) ForgetPassword(ctx context.Context, email string) (interface{}, error) {
	var user *model.User
	if err := r.DB.First(&user, "email = ?", email).Error; err != nil {
		return map[string]interface{}{
			"status": "email not found",
		}, err
	}

	reset := &model.Reset{
		ID:     uuid.NewString(),
		UserId: user.ID,
	}
	r.resets = append(r.resets, reset)
	r.DB.Create(reset)

	SendResetPasswordLink(r.Resolver, user, reset)

	return map[string]interface{}{
		"status": "email has been sent",
	}, nil
}

// ResetPassword is the resolver for the ResetPassword field.
func (r *mutationResolver) ResetPassword(ctx context.Context, id string, password string) (interface{}, error) {
	var reset *model.Reset
	if err := r.DB.First(&reset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user, err := UserById(r.Resolver, reset.UserId)
	if err != nil {
		return nil, err
	}
	user.Password = password
	r.DB.Save(user)
	r.DB.Delete(reset)
	return map[string]interface{}{
		"status": "success",
	}, nil
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
	follow := &model.UserFollow{
		ID:       uuid.NewString(),
		UserId:   id1,
		FollowId: id2,
	}
	r.user_follows = append(r.user_follows, follow)
	r.DB.Create(follow)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// UnFollow is the resolver for the UnFollow field.
func (r *mutationResolver) UnFollow(ctx context.Context, id1 string, id2 string) (interface{}, error) {
	var follow *model.UserFollow
	if err := r.DB.First(&follow, "user_id = ? AND follow_id = ?", id1, id2).Error; err != nil {
		return map[string]interface{}{
			"status": "not found",
		}, nil
	}
	r.user_follows = lo.Filter[*model.UserFollow](r.user_follows, func(x *model.UserFollow, _ int) bool {
		return x.UserId == id1 && x.FollowId == id2
	})
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

// SendMessage is the resolver for the SendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.InputMessage) (interface{}, error) {
	message := &model.Message{
		ID:        uuid.NewString(),
		Text:      input.Text,
		User1Id:   input.User1Id,
		User2Id:   input.User2Id,
		CreatedAt: time.Time{},
	}
	r.messages = append(r.messages, message)
	r.DB.Create(message)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// Visit is the resolver for the Visit field.
func (r *mutationResolver) Visit(ctx context.Context, id string) (interface{}, error) {
	myId := auth.JwtGetValue(ctx).Userid

	var visit *model.UserVisit
	err := r.DB.First(&visit, "visit_id = ? and user_id = ?", id, myId).Error
	if err == nil {

		var visits []*model.UserVisit
		r.DB.Find(&visits, "visit_id = ?", id)
		return map[string]interface{}{
			"status": "already exists",
			"length": len(visits),
		}, nil
	}
	visit = &model.UserVisit{
		ID:      uuid.NewString(),
		UserId:  myId,
		VisitId: id,
	}
	r.user_visits = append(r.user_visits, visit)
	r.DB.Create(visit)

	var visits []*model.UserVisit
	r.DB.Find(&visits, "visit_id = ?", id)

	return map[string]interface{}{
		"status": "success",
		"length": len(visits),
	}, nil
}

// AddEducation is the resolver for the AddEducation field.
func (r *mutationResolver) AddEducation(ctx context.Context, input model.InputEducation) (*model.Education, error) {
	education := &model.Education{
		ID:        uuid.NewString(),
		School:    input.School,
		Field:     input.Field,
		StartedAt: input.StartedAt,
		EndedAt:   input.EndedAt,
	}
	r.educations = append(r.educations, education)
	r.DB.Create(education)

	userEducation := &model.UserEducation{
		ID:          uuid.NewString(),
		UserId:      auth.JwtGetValue(ctx).Userid,
		EducationId: education.ID,
	}
	r.user_educations = append(r.user_educations, userEducation)
	r.DB.Create(userEducation)

	return education, nil
}

// UpdateEducation is the resolver for the UpdateEducation field.
func (r *mutationResolver) UpdateEducation(ctx context.Context, id string, input model.InputEducation) (*model.Education, error) {
	var education *model.Education
	if err := r.DB.First(&education, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if input.School != "" {
		education.School = input.School
	}
	if input.Field != "" {
		education.Field = input.Field
	}
	if input.StartedAt != "" {
		education.StartedAt = input.StartedAt
	}
	if input.EndedAt != "" {
		education.EndedAt = input.EndedAt
	}

	r.DB.Save(education)

	return education, nil
}

// RemoveEducation is the resolver for the RemoveEducation field.
func (r *mutationResolver) RemoveEducation(ctx context.Context, id string) (interface{}, error) {
	var education *model.Education
	if err := r.DB.First(&education, "id = ?", id).Error; err != nil {
		return nil, err
	}
	r.DB.Delete(education)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// AddExperience is the resolver for the AddExperience field.
func (r *mutationResolver) AddExperience(ctx context.Context, input model.InputExperience) (*model.Experience, error) {
	experience := &model.Experience{
		ID:        uuid.NewString(),
		Position:  input.Position,
		Desc:      input.Desc,
		Company:   input.Company,
		StartedAt: input.StartedAt,
		EndedAt:   input.EndedAt,
		IsActive:  input.IsActive,
	}
	r.experiences = append(r.experiences, experience)
	r.DB.Create(experience)

	userExperience := &model.UserExperience{
		ID:           uuid.NewString(),
		UserId:       auth.JwtGetValue(ctx).Userid,
		ExperienceId: experience.ID,
	}
	r.user_experiences = append(r.user_experiences, userExperience)
	r.DB.Create(userExperience)

	return experience, nil
}

// UpdateExperience is the resolver for the UpdateExperience field.
func (r *mutationResolver) UpdateExperience(ctx context.Context, id string, input model.InputExperience) (*model.Experience, error) {
	var experience *model.Experience
	if err := r.DB.First(&experience, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if input.StartedAt != "" {
		experience.StartedAt = input.StartedAt
	}
	if input.EndedAt != "" {
		experience.EndedAt = input.EndedAt
	}
	if input.Position != "" {
		experience.Position = input.Position
	}
	if input.Desc != "" {
		experience.Desc = input.Desc
	}
	if input.Company != "" {
		experience.Company = input.Company
	}
	if input.IsActive != experience.IsActive {
		experience.IsActive = input.IsActive
	}

	r.DB.Save(experience)

	return experience, nil
}

// RemoveExperience is the resolver for the RemoveExperience field.
func (r *mutationResolver) RemoveExperience(ctx context.Context, id string) (interface{}, error) {
	var experience *model.Experience
	if err := r.DB.First(&experience, "id = ?", id).Error; err != nil {
		return nil, err
	}
	fmt.Println(experience)
	r.DB.Delete(experience)

	return map[string]interface{}{
		"status": "success",
	}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return UserById(r.Resolver, id)
}

// UsersByName is the resolver for the UsersByName field.
func (r *queryResolver) UsersByName(ctx context.Context, name *string, limit int, offset int) ([]*model.User, error) {
	if *name == "" {
		*name = "%"
	}

	var users []*model.User
	if err := r.DB.Limit(limit).Offset(offset).Find(&users, "concat(first_name, mid_name, last_name) like ?", "%"+*name+"%").Error; err != nil {
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
		Visits:          []*model.User{},
		Follows:         []*model.User{},
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

// CheckReset is the resolver for the CheckReset field.
func (r *queryResolver) CheckReset(ctx context.Context, id string) (*model.User, error) {
	var reset *model.Reset
	if err := r.DB.First(&reset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	user, err := UserById(r.Resolver, reset.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// IsFollow is the resolver for the IsFollow field.
func (r *queryResolver) IsFollow(ctx context.Context, id1 string, id2 string) (bool, error) {
	var follow *model.UserFollow
	if err := r.DB.First(&follow, "user_id = ? AND follow_id = ?", id1, id2).Error; err != nil {
		return false, err
	}
	if follow != nil {
		return true, nil
	}
	return false, nil
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

// ConnectedUsers is the resolver for the ConnectedUsers field.
func (r *queryResolver) ConnectedUsers(ctx context.Context) ([]*model.User, error) {
	myId := auth.JwtGetValue(ctx).Userid

	var userIds []string
	var err error
	var connections []*model.Connection
	err = r.DB.Find(&connections, "user1_id = ?", myId).Error
	if err != nil {
		return nil, err
	}
	userIds = lo.Map[*model.Connection, string](connections, func(x *model.Connection, _ int) string {
		return x.User2ID
	})

	err = r.DB.Find(&connections, "user2_id = ?", myId).Error
	if err != nil {
		return nil, err
	}
	userIds = append(userIds, lo.Map[*model.Connection, string](connections, func(x *model.Connection, _ int) string {
		return x.User1ID
	})...)

	return UsersById(r.Resolver, userIds)
}

// Messages is the resolver for the Messages field.
func (r *queryResolver) Messages(ctx context.Context, id1 string, id2 string) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.DB.Find(&messages, "(user1_id = ? and user2_id = ?) or (user1_id = ? and user2_id = ?) order by created_at", id1, id2, id2, id1).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// User is the resolver for the User field.
func (r *resetResolver) User(ctx context.Context, obj *model.Reset) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Visits is the resolver for the Visits field.
func (r *userResolver) Visits(ctx context.Context, obj *model.User) ([]*model.User, error) {
	var visits []*model.UserVisit
	err := r.DB.Find(&visits, "visit_id = ?", obj.ID).Error
	if err != nil {
		return nil, err
	}

	visitIds := lo.Map[*model.UserVisit, string](visits, func(x *model.UserVisit, _ int) string {
		return x.UserId
	})
	return UsersById(r.Resolver, visitIds)
}

// Follows is the resolver for the Follows field.
func (r *userResolver) Follows(ctx context.Context, obj *model.User) ([]*model.User, error) {
	var follows []*model.UserFollow
	err := r.DB.Find(&follows, "user_id = ?", obj.ID).Error
	if err != nil {
		return nil, err
	}

	followIds := lo.Map[*model.UserFollow, string](follows, func(x *model.UserFollow, _ int) string {
		return x.FollowId
	})
	return UsersById(r.Resolver, followIds)
}

// Experiences is the resolver for the Experiences field.
func (r *userResolver) Experiences(ctx context.Context, obj *model.User) ([]*model.Experience, error) {
	experienceIds := lo.Map[*model.Experience, string](obj.Experiences, func(x *model.Experience, _ int) string {
		return x.ID
	})

	var experiences []*model.Experience
	if err := r.DB.Find(&experiences, experienceIds).Error; err != nil {
		return nil, err
	}
	return experiences, nil
}

// Educations is the resolver for the Educations field.
func (r *userResolver) Educations(ctx context.Context, obj *model.User) ([]*model.Education, error) {
	educationIds := lo.Map[*model.Education, string](obj.Educations, func(x *model.Education, _ int) string {
		return x.ID
	})

	var educations []*model.Education
	if err := r.DB.Find(&educations, educationIds).Error; err != nil {
		return nil, err
	}
	return educations, nil
}

// Activation returns generated.ActivationResolver implementation.
func (r *Resolver) Activation() generated.ActivationResolver { return &activationResolver{r} }

// ConnectRequest returns generated.ConnectRequestResolver implementation.
func (r *Resolver) ConnectRequest() generated.ConnectRequestResolver {
	return &connectRequestResolver{r}
}

// Connection returns generated.ConnectionResolver implementation.
func (r *Resolver) Connection() generated.ConnectionResolver { return &connectionResolver{r} }

// Message returns generated.MessageResolver implementation.
func (r *Resolver) Message() generated.MessageResolver { return &messageResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Reset returns generated.ResetResolver implementation.
func (r *Resolver) Reset() generated.ResetResolver { return &resetResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type activationResolver struct{ *Resolver }
type connectRequestResolver struct{ *Resolver }
type connectionResolver struct{ *Resolver }
type messageResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resetResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
