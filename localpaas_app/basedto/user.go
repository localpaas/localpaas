package basedto

import "github.com/localpaas/localpaas/localpaas_app/entity"

type UserBaseResp struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Photo    string `json:"photo"`
}

func TransformUserBase(user *entity.User) *UserBaseResp {
	if user == nil {
		return nil
	}
	return &UserBaseResp{
		ID:       user.ID,
		FullName: user.FullName,
		Photo:    user.Photo,
	}
}

func TransformUsersBase(users []*entity.User) []*UserBaseResp {
	resp, _ := TransformObjectSlice(users, func(user *entity.User) (*UserBaseResp, error) {
		return TransformUserBase(user), nil
	})
	return resp
}
