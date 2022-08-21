package graph

import "MZ221-TPA-Web-Back/graph/model"

func UserById(r *Resolver, id string) (*model.User, error) {
	var user *model.User
	if err := r.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UsersById(r *Resolver, ids []string) ([]*model.User, error) {
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
