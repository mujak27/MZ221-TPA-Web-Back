package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"MZ221-TPA-Web-Back/auth"
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"errors"
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

// User is the resolver for the User field.
func (r *activationResolver) User(ctx context.Context, obj *model.Activation) (*model.User, error) {
	return UserById(r.Resolver, obj.UserId)
}

// User is the resolver for the User field.
func (r *activityResolver) User(ctx context.Context, obj *model.Activity) (*model.User, error) {
	return UserById(r.Resolver, obj.UserId)
}

// User1 is the resolver for the User1 field.
func (r *blockResolver) User1(ctx context.Context, obj *model.Block) (*model.User, error) {
	return UserById(r.Resolver, obj.User1.ID)
}

// User2 is the resolver for the User2 field.
func (r *blockResolver) User2(ctx context.Context, obj *model.Block) (*model.User, error) {
	return UserById(r.Resolver, obj.User2.ID)
}

// User1 is the resolver for the User1 field.
func (r *connectRequestResolver) User1(ctx context.Context, obj *model.ConnectRequest) (*model.User, error) {
	return UserById(r.Resolver, obj.User1ID)
}

// User2 is the resolver for the User2 field.
func (r *connectRequestResolver) User2(ctx context.Context, obj *model.ConnectRequest) (*model.User, error) {
	return UserById(r.Resolver, obj.User2ID)
}

// User1 is the resolver for the User1 field.
func (r *connectionResolver) User1(ctx context.Context, obj *model.Connection) (*model.User, error) {
	return UserById(r.Resolver, obj.User1ID)
}

// User2 is the resolver for the User2 field.
func (r *connectionResolver) User2(ctx context.Context, obj *model.Connection) (*model.User, error) {
	return UserById(r.Resolver, obj.User2ID)
}

// User is the resolver for the User field.
func (r *jobResolver) User(ctx context.Context, obj *model.Job) (*model.User, error) {
	return UserById(r.Resolver, obj.UserId)
}

// LoginRegisWithSso is the resolver for the LoginRegisWithSSO field.
func (r *mutationResolver) LoginRegisWithSso(ctx context.Context, googleToken string) (string, error) {
	tokenValue, _ := jwt.Parse(googleToken, nil)

	// validate the essential claims
	// if err != nil {
	// 	return "", err
	// }

	email := tokenValue.Claims.(jwt.MapClaims)["email"].(string)
	firstName := tokenValue.Claims.(jwt.MapClaims)["given_name"].(string)
	lastName := tokenValue.Claims.(jwt.MapClaims)["family_name"].(string)
	isActive := tokenValue.Claims.(jwt.MapClaims)["email_verified"].(bool)
	profilePhoto := tokenValue.Claims.(jwt.MapClaims)["picture"].(string)
	fmt.Println(email)
	fmt.Println(firstName)
	fmt.Println(lastName)
	fmt.Println(isActive)
	fmt.Println(profilePhoto)
	if email == "" {
		return "", errors.New("invalid email")
	}

	var user *model.User
	err := r.DB.First(&user, "email = ?", email).Error
	log.Println(user)
	log.Println(err)
	if user.ID == "" {
		fmt.Println("not found")
		user = &model.User{
			ID:       uuid.NewString(),
			Email:    email,
			Password: "",

			FirstName:       firstName,
			LastName:        lastName,
			MidName:         "",
			IsActive:        isActive,
			ProfilePhoto:    profilePhoto,
			BackgroundPhoto: "",
			Headline:        "",
			Pronoun:         "",
			ProfileLink:     "",
			About:           "",
			Location:        "",
			IsSso:           true,
			HasFilledData:   false,
		}
		user.ProfileLink = user.ID

		fmt.Println(user)

		r.users = append(r.users, user)

		r.DB.Create(user)

		if !user.IsActive {
			SendActivationLink(r.Resolver, user)
		}
	}

	fmt.Println("done")

	token, err := auth.JwtGenerate(ctx, user.ID)
	if err != nil {
		return "", err
	}

	fmt.Println(token)

	return token, nil
}

// Login is the resolver for the Login field.
func (r *mutationResolver) Login(ctx context.Context, input model.InputLogin) (string, error) {
	var user *model.User
	if err := r.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		return "", err
	}

	if user.IsSso {
		return "", errors.New("set password via sso login first")
	}

	token, err := auth.JwtGenerate(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Register is the resolver for the Register field.
func (r *mutationResolver) Register(ctx context.Context, input *model.InputRegister) (string, error) {
	if !validateEmail(input.Email) {
		return "", errors.New("invalid email")
	}

	if len(input.Password) < 8 {
		return "", errors.New("minimum length for password is 8 characters")
	}

	var user *model.User
	if err := r.DB.First(&user, "email = ?", input.Email).Error; err == nil {
		return "", errors.New("email already registered")
	}

	user = &model.User{
		ID:              uuid.NewString(),
		Email:           input.Email,
		Password:        input.Password,
		FirstName:       "",
		LastName:        "",
		MidName:         "",
		IsActive:        false,
		ProfilePhoto:    "",
		BackgroundPhoto: "",
		Headline:        "",
		Pronoun:         "",
		ProfileLink:     "",
		About:           "",
		Location:        "",
		IsSso:           false,
		HasFilledData:   false,
	}
	user.ProfileLink = user.ID
	r.users = append(r.users, user)

	r.DB.Create(user)

	SendActivationLink(r.Resolver, user)

	token, err := auth.JwtGenerate(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// FirstUpdateProfile is the resolver for the FirstUpdateProfile field.
func (r *mutationResolver) FirstUpdateProfile(ctx context.Context, input model.InputFirstUpdateProfile) (model.MutationStatus, error) {
	myId := getId(ctx)

	user, err := UserById(r.Resolver, myId)
	if err != nil {
		return model.MutationStatusNotFound, err
	}

	fmt.Println((user))

	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.MidName = input.MidName
	user.ProfilePhoto = input.ProfilePhoto
	user.Pronoun = input.Pronoun
	user.HasFilledData = true

	r.DB.Save(user)

	return model.MutationStatusSuccess, nil
}

// UpdateProfile is the resolver for the UpdateProfile field.
func (r *mutationResolver) UpdateProfile(ctx context.Context, input model.InputUser) (model.MutationStatus, error) {
	myId := getId(ctx)

	user, err := UserByProfileLink(ctx, r.Resolver, input.ProfileLink)
	if err == nil && user.ID != myId {
		return model.MutationStatusError, errors.New("profile link already taken")
	}
	user, err = UserById(r.Resolver, myId)
	if err != nil || user == nil {
		return model.MutationStatusNotFound, err
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
	user.ProfileLink = input.ProfileLink

	r.DB.Save(user)

	return model.MutationStatusSuccess, nil
}

// ForgetPassword is the resolver for the ForgetPassword field.
func (r *mutationResolver) ForgetPassword(ctx context.Context, email string) (model.MutationStatus, error) {
	var user *model.User
	if err := r.DB.First(&user, "email = ?", email).Error; err != nil {
		return model.MutationStatusNotFound, errors.New(string(model.MutationStatusNotFound))
	}

	reset := &model.Reset{
		ID:     uuid.NewString(),
		UserId: user.ID,
	}
	r.resets = append(r.resets, reset)
	r.DB.Create(reset)

	SendResetPasswordLink(r.Resolver, user, reset)

	return model.MutationStatusSuccess, nil
}

// ResetPassword is the resolver for the ResetPassword field.
func (r *mutationResolver) ResetPassword(ctx context.Context, id string, password string) (model.MutationStatus, error) {
	var reset *model.Reset
	if err := r.DB.First(&reset, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
	}
	user, err := UserById(r.Resolver, reset.UserId)
	if err != nil {
		return model.MutationStatusNotFound, err
	}
	user.Password = password
	user.IsSso = false
	r.DB.Save(user)
	r.DB.Delete(reset)
	return model.MutationStatusSuccess, nil
}

// FirstFillData is the resolver for the FirstFillData field.
func (r *mutationResolver) FirstFillData(ctx context.Context, input model.InputUser) (model.MutationStatus, error) {
	panic(fmt.Errorf("not implemented"))
}

// ChangePassword is the resolver for the ChangePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, password string) (model.MutationStatus, error) {
	myId := getId(ctx)

	user, _ := UserById(r.Resolver, myId)

	user.IsSso = false
	user.Password = password

	r.DB.Save(user)

	return model.MutationStatusSuccess, nil
}

// SendActivation is the resolver for the SendActivation field.
func (r *mutationResolver) SendActivation(ctx context.Context, id string) (model.MutationStatus, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
	}
	SendActivationLink(r.Resolver, user)
	return model.MutationStatusSuccess, nil
}

// Activate is the resolver for the Activate field.
func (r *mutationResolver) Activate(ctx context.Context, id string) (model.MutationStatus, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
	}
	user.IsActive = true
	r.DB.Save(user)
	return model.MutationStatusSuccess, nil
}

// Block is the resolver for the Block field.
func (r *mutationResolver) Block(ctx context.Context, userID string) (*model.User, error) {
	myId := getId(ctx)
	user, _ := UserById(r.Resolver, userID)
	block := &model.Block{
		ID:      uuid.NewString(),
		User1Id: myId,
		User2Id: userID,
	}
	r.blocks = append(r.blocks, block)
	r.DB.Create(block)
	return user, nil
}

// UnBlock is the resolver for the UnBlock field.
func (r *mutationResolver) UnBlock(ctx context.Context, userID string) (*model.User, error) {
	myId := getId(ctx)
	user, _ := UserById(r.Resolver, userID)
	var block *model.Block
	if err := r.DB.First(&block, "user1_id = ? AND user2_id = ?", myId, userID).Error; err != nil {
		return nil, err
	}
	r.blocks = lo.Filter[*model.Block](r.blocks, func(x *model.Block, _ int) bool {
		return x.User1Id == myId && x.User2Id == userID
	})
	r.DB.Delete(block)
	return user, nil
}

// Follow is the resolver for the Follow field.
func (r *mutationResolver) Follow(ctx context.Context, id1 string, id2 string) (model.MutationStatus, error) {
	follow := &model.UserFollow{
		ID:       uuid.NewString(),
		UserId:   id1,
		FollowId: id2,
	}
	r.user_follows = append(r.user_follows, follow)
	r.DB.Create(follow)

	myUser, _ := UserById(r.Resolver, id1)

	activityText := myUser.FirstName + " " + myUser.LastName + " has followed you"
	AddActivity(r.Resolver, id2, activityText)
	return model.MutationStatusSuccess, nil
}

// UnFollow is the resolver for the UnFollow field.
func (r *mutationResolver) UnFollow(ctx context.Context, id1 string, id2 string) (model.MutationStatus, error) {
	var follow *model.UserFollow
	if err := r.DB.First(&follow, "user_id = ? AND follow_id = ?", id1, id2).Error; err != nil {
		return model.MutationStatusNotFound, nil
	}
	r.user_follows = lo.Filter[*model.UserFollow](r.user_follows, func(x *model.UserFollow, _ int) bool {
		return x.UserId == id1 && x.FollowId == id2
	})
	r.DB.Delete(follow)
	return model.MutationStatusSuccess, nil
}

// SendConnectRequest is the resolver for the SendConnectRequest field.
func (r *mutationResolver) SendConnectRequest(ctx context.Context, id1 string, id2 string) (model.MutationStatus, error) {
	var connectRequest *model.ConnectRequest
	err := r.DB.First(&connectRequest, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err == nil {
		return model.MutationStatusAlreadyExist, nil
	}
	connectRequest = &model.ConnectRequest{
		ID:      uuid.NewString(),
		User1ID: id1,
		User2ID: id2,
	}

	r.connectRequests = append(r.connectRequests, connectRequest)

	r.DB.Create(connectRequest)

	return model.MutationStatusSuccess, nil
}

// DeleteConnectRequest is the resolver for the DeleteConnectRequest field.
func (r *mutationResolver) DeleteConnectRequest(ctx context.Context, id1 string, id2 string) (model.MutationStatus, error) {
	var connectRequest *[]model.ConnectRequest
	err1 := r.DB.Find(&connectRequest, "user1_id = ? and user2_id = ?", id1, id2).Error
	r.DB.Delete(connectRequest)
	err2 := r.DB.Find(&connectRequest, "user1_id = ? and user2_id = ?", id2, id1).Error
	r.DB.Delete(connectRequest)
	if err1 != nil && err2 != nil {
		return model.MutationStatusNotFound, err1
	}
	// r.DB.Delete(connectRequest)
	return model.MutationStatusSuccess, nil
}

// AcceptConnectRequest is the resolver for the AcceptConnectRequest field.
func (r *mutationResolver) AcceptConnectRequest(ctx context.Context, id1 string, id2 string) (model.MutationStatus, error) {
	_, err := r.DeleteConnectRequest(ctx, id1, id2)
	if err != nil {
		return model.MutationStatusNotFound, err
	}

	id1, id2 = SortIdAsc(id1, id2)

	connection := &model.Connection{
		ID:      uuid.NewString(),
		User1ID: id1,
		User2ID: id2,
	}
	r.connections = append(r.connections, connection)
	r.DB.Create(connection)

	return model.MutationStatusSuccess, nil
}

// UnConnect is the resolver for the UnConnect field.
func (r *mutationResolver) UnConnect(ctx context.Context, id1 string, id2 string) (model.MutationStatus, error) {
	id1, id2 = SortIdAsc(id1, id2)

	var connection *model.Connection
	err := r.DB.First(&connection, "user1_id = ? and user2_id = ?", id1, id2).Error
	if err != nil {
		return model.MutationStatusNotFound, err
	}
	r.DB.Delete(connection)
	return model.MutationStatusSuccess, nil
}

// Visit is the resolver for the Visit field.
func (r *mutationResolver) Visit(ctx context.Context, id string) (model.MutationStatus, error) {
	myId := auth.JwtGetValue(ctx).Userid

	var visit *model.UserVisit
	err := r.DB.First(&visit, "visit_id = ? and user_id = ?", id, myId).Error
	if err == nil {
		var visits []*model.UserVisit
		r.DB.Find(&visits, "visit_id = ?", id)
		return model.MutationStatusAlreadyExist, nil
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

	return model.MutationStatusSuccess, nil
}

// VisitByLink is the resolver for the VisitByLink field.
func (r *mutationResolver) VisitByLink(ctx context.Context, profileLink string) (model.MutationStatus, error) {
	myId := auth.JwtGetValue(ctx).Userid

	user, err := UserByProfileLink(ctx, r.Resolver, profileLink)
	if err != nil {
		return model.MutationStatusError, errors.New("not found")
	}
	id := user.ID

	var visit *model.UserVisit
	err = r.DB.First(&visit, "visit_id = ? and user_id = ?", id, myId).Error
	if err == nil {
		var visits []*model.UserVisit
		r.DB.Find(&visits, "visit_id = ?", id)
		return model.MutationStatusAlreadyExist, nil
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

	return model.MutationStatusSuccess, nil
}

// AddEducation is the resolver for the AddEducation field.
func (r *mutationResolver) AddEducation(ctx context.Context, input model.InputEducation) (model.MutationStatus, error) {
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

	return model.MutationStatusSuccess, nil
}

// UpdateEducation is the resolver for the UpdateEducation field.
func (r *mutationResolver) UpdateEducation(ctx context.Context, id string, input model.InputEducation) (model.MutationStatus, error) {
	var education *model.Education
	if err := r.DB.First(&education, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
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

	return model.MutationStatusSuccess, nil
}

// RemoveEducation is the resolver for the RemoveEducation field.
func (r *mutationResolver) RemoveEducation(ctx context.Context, id string) (model.MutationStatus, error) {
	var education *model.Education
	if err := r.DB.First(&education, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
	}
	r.DB.Delete(education)

	return model.MutationStatusSuccess, nil
}

// AddExperience is the resolver for the AddExperience field.
func (r *mutationResolver) AddExperience(ctx context.Context, input model.InputExperience) (model.MutationStatus, error) {
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

	return model.MutationStatusSuccess, nil
}

// UpdateExperience is the resolver for the UpdateExperience field.
func (r *mutationResolver) UpdateExperience(ctx context.Context, id string, input model.InputExperience) (model.MutationStatus, error) {
	var experience *model.Experience
	if err := r.DB.First(&experience, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
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

	return model.MutationStatusSuccess, nil
}

// RemoveExperience is the resolver for the RemoveExperience field.
func (r *mutationResolver) RemoveExperience(ctx context.Context, id string) (model.MutationStatus, error) {
	var experience *model.Experience
	if err := r.DB.First(&experience, "id = ?", id).Error; err != nil {
		return model.MutationStatusNotFound, err
	}
	fmt.Println(experience)
	r.DB.Delete(experience)

	return model.MutationStatusSuccess, nil
}

// AddJob is the resolver for the AddJob field.
func (r *mutationResolver) AddJob(ctx context.Context, text string) (model.MutationStatus, error) {
	myId := getId(ctx)
	job := &model.Job{
		ID:     uuid.NewString(),
		UserId: myId,
		Text:   text,
	}
	r.jobs = append(r.jobs, job)
	r.DB.Create(job)

	user, _ := UserById(r.Resolver, myId)

	activityText := user.FirstName + " " + user.LastName + " has posted a new job : " + text
	AddActivity(r.Resolver, myId, activityText)

	return model.MutationStatusSuccess, nil
}

// CountUser is the resolver for the CountUser field.
func (r *queryResolver) CountUser(ctx context.Context, keyword *string) (int, error) {
	users, nil := r.UsersByName(ctx, keyword, 1000000000, 0)
	return len(users), nil
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

	fmt.Println(*name)
	var users []*model.User
	if err := r.DB.Limit(limit).Offset(offset).Find(&users, "concat(first_name,mid_name,last_name) like ?", "%"+*name+"%").Error; err != nil {
		return nil, err
	}
	fmt.Println(*name)
	fmt.Println(users)
	return users, nil
}

// UserByLink is the resolver for the UserByLink field.
func (r *queryResolver) UserByLink(ctx context.Context, link string) (*model.User, error) {
	return UserByProfileLink(ctx, r.Resolver, link)
}

// IsEmailValid is the resolver for the isEmailValid field.
func (r *queryResolver) IsEmailValid(ctx context.Context, email string) (bool, error) {
	var user *model.User
	if err := r.DB.First(&user, "email = ?", email).Error; err != nil {
		return false, err
	}
	return true, nil
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

// IsBlock is the resolver for the IsBlock field.
func (r *queryResolver) IsBlock(ctx context.Context, userID string) (bool, error) {
	myId := getId(ctx)
	var block *model.Block
	if err := r.DB.First(&block, "user1_id = ? and user2_id = ?", myId, userID).Error; err != nil {
		return false, err
	}
	return true, nil
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

// ConnectionRequest is the resolver for the ConnectionRequest field.
func (r *queryResolver) ConnectionRequest(ctx context.Context) ([]*model.User, error) {
	myId := getId(ctx)

	var connectRequests []*model.ConnectRequest
	err := r.DB.Find(&connectRequests, "user2_id = ?", myId).Error
	if err != nil {
		return nil, err
	}
	fmt.Println(connectRequests)

	userIds := lo.Map[*model.ConnectRequest, string](connectRequests, func(x *model.ConnectRequest, _ int) string {
		return x.User1ID
	})
	fmt.Println(userIds)
	return UsersById(r.Resolver, userIds)
}

// ConnectedUsers is the resolver for the ConnectedUsers field.
func (r *queryResolver) ConnectedUsers(ctx context.Context) ([]*model.User, error) {
	userIds, err := getConnectedIds(r.Resolver, ctx)
	if err != nil {
		return nil, err
	}

	return UsersById(r.Resolver, userIds)
}

// Activities is the resolver for the Activities field.
func (r *queryResolver) Activities(ctx context.Context) ([]*model.Activity, error) {
	userIds, err := getFriendIds(r.Resolver, ctx)
	if err != nil {
		return nil, err
	}

	var activities []*model.Activity
	if err := r.DB.Find(&activities, "user_id IN ?", userIds).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

// UsersSuggestion is the resolver for the UsersSuggestion field.
func (r *queryResolver) UsersSuggestion(ctx context.Context) ([]*model.User, error) {
	myId := getId(ctx)

	var connections []*model.Connection
	if err := r.DB.Find(&connections, "user1_id = ? or user2_id = ?", myId, myId).Error; err != nil {
		return nil, err
	}
	connectedUserIds := lo.Map[*model.Connection, []string](connections, func(x *model.Connection, _ int) []string {
		return []string{x.User2ID, x.User1ID}
	})
	flattenConnectedUserIds := lo.Flatten[string](connectedUserIds)
	flattenConnectedUserIds = lo.Uniq[string](flattenConnectedUserIds)

	if err := r.DB.Find(&connections, "user1_id IN ? or user2_id IN ?", flattenConnectedUserIds, flattenConnectedUserIds).Error; err != nil {
		return nil, err
	}
	suggestedUserIds := lo.Map[*model.Connection, []string](connections, func(x *model.Connection, _ int) []string {
		return []string{x.User2ID, x.User1ID}
	})
	flattenSuggestedUserIds := lo.Flatten[string](suggestedUserIds)
	flattenSuggestedUserIds = lo.Uniq[string](flattenSuggestedUserIds)

	flattenSuggestedUserIds = lo.Filter[string](flattenSuggestedUserIds, func(x string, _ int) bool {
		return !lo.Contains[string](append(flattenConnectedUserIds, myId), x)
	})

	return UsersById(r.Resolver, flattenSuggestedUserIds)
}

// Jobs is the resolver for the Jobs field.
func (r *queryResolver) Jobs(ctx context.Context) ([]*model.Job, error) {
	var jobs []*model.Job
	if err := r.DB.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// User is the resolver for the User field.
func (r *resetResolver) User(ctx context.Context, obj *model.Reset) (*model.User, error) {
	return UserById(r.Resolver, obj.UserId)
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
	fmt.Println(obj.ID)
	var userExperiences []*model.UserExperience
	if err := r.DB.Find(&userExperiences, "user_id = ?", obj.ID).Error; err != nil {
		return nil, err
	}
	fmt.Println(userExperiences)

	experienceIds := lo.Map[*model.UserExperience, string](userExperiences, func(x *model.UserExperience, _ int) string {
		return x.ExperienceId
	})

	var experiences []*model.Experience
	if err := r.DB.Find(&experiences, "id in ?", experienceIds).Error; err != nil {
		return nil, err
	}
	return experiences, nil
}

// Educations is the resolver for the Educations field.
func (r *userResolver) Educations(ctx context.Context, obj *model.User) ([]*model.Education, error) {
	var userEducations []*model.UserEducation
	if err := r.DB.Find(&userEducations, "user_id = ?", obj.ID).Error; err != nil {
		return nil, err
	}
	fmt.Println(userEducations)

	educationIds := lo.Map[*model.UserEducation, string](userEducations, func(x *model.UserEducation, _ int) string {
		return x.EducationId
	})

	var educations []*model.Education
	if err := r.DB.Find(&educations, "id in ?", educationIds).Error; err != nil {
		return nil, err
	}
	return educations, nil
}

// Activation returns generated.ActivationResolver implementation.
func (r *Resolver) Activation() generated.ActivationResolver { return &activationResolver{r} }

// Activity returns generated.ActivityResolver implementation.
func (r *Resolver) Activity() generated.ActivityResolver { return &activityResolver{r} }

// Block returns generated.BlockResolver implementation.
func (r *Resolver) Block() generated.BlockResolver { return &blockResolver{r} }

// ConnectRequest returns generated.ConnectRequestResolver implementation.
func (r *Resolver) ConnectRequest() generated.ConnectRequestResolver {
	return &connectRequestResolver{r}
}

// Connection returns generated.ConnectionResolver implementation.
func (r *Resolver) Connection() generated.ConnectionResolver { return &connectionResolver{r} }

// Job returns generated.JobResolver implementation.
func (r *Resolver) Job() generated.JobResolver { return &jobResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Reset returns generated.ResetResolver implementation.
func (r *Resolver) Reset() generated.ResetResolver { return &resetResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type activationResolver struct{ *Resolver }
type activityResolver struct{ *Resolver }
type blockResolver struct{ *Resolver }
type connectRequestResolver struct{ *Resolver }
type connectionResolver struct{ *Resolver }
type jobResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type resetResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
