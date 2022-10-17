package graph

import (
	"MZ221-TPA-Web-Back/auth"
	"MZ221-TPA-Web-Back/graph/model"
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

func AddActivity(r *Resolver, userId string, text string) error {
	activity := &model.Activity{
		ID:     uuid.NewString(),
		UserId: userId,
		Text:   text,
	}
	r.activities = append(r.activities, activity)
	return r.DB.Create(activity).Error
}

func UserByProfileLink(ctx context.Context, r *Resolver, link string) (*model.User, error) {
	myId := getId(ctx)
	var user *model.User
	if err := r.DB.First(&user, "profile_link = ?", link).Error; err != nil {
		return nil, err
	}
	var block *model.Block
	if err := r.DB.First(&block, "user1_id = ? and user2_id = ?", user.ID, myId).Error; err == nil {
		return nil, err
	}
	return user, nil
}

func UserById(r *Resolver, id string) (*model.User, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UsersById(r *Resolver, ids []string) ([]*model.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var users []*model.User
	if err := r.DB.Find(&users, ids).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func SortIdAsc(id1 string, id2 string) (string, string) {
	if id1 < id2 {
		return id1, id2
	}
	return id2, id1
}

func getId(ctx context.Context) string {
	return auth.JwtGetValue(ctx).Userid

}

func getBlockIds(r *Resolver, ctx context.Context) ([]string, error) {

	myId := getId(ctx)
	var blocks []*model.Block
	if err := r.DB.Select("user1_id").Find(&blocks, "user2_id = ?", myId).Error; err != nil {
		return nil, err
	}
	return lo.Map[*model.Block, string](blocks, func(x *model.Block, _ int) string {
		return x.User1Id
	}), nil
}

func getConnectedIds(r *Resolver, ctx context.Context) ([]string, error) {

	var idList []string
	myId := auth.JwtGetValue(ctx).Userid

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

	return idList, nil
}

func getFriendIds(r *Resolver, ctx context.Context) ([]string, error) {

	var idList []string
	myId := auth.JwtGetValue(ctx).Userid
	idList = append(idList, myId)

	user, err := UserById(r, myId)
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

	return idList, err
}

func validateEmail(email string) bool {

	indexAt := strings.Index(email, "@")
	indexDot := strings.Index(email, ".")
	length := len(email)

	if length == 0 || indexAt == -1 || indexDot == -1 || indexAt > indexDot || indexAt == indexDot-1 || indexAt == 0 || indexDot == length-1 {
		return false
	}

	return true
}

func concatUserName(user *model.User) string {
	res := user.FirstName
	if user.MidName != "" {
		res += " " + user.MidName
	}
	if user.LastName != "" {
		res += " " + user.LastName
	}

	return res
}
