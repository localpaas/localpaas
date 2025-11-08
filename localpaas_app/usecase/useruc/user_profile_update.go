package useruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *UserUC) UpdateProfile(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.UpdateProfileReq,
) (*userdto.UpdateProfileResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		profileData := &userProfileData{}
		err := uc.loadUserProfileData(ctx, db, auth, req, profileData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingUserProfileData{}
		uc.preparePersistingUserProfileData(req, profileData, persistingData)

		return uc.userRepo.Update(ctx, db, persistingData.UpdatingUser)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.UpdateProfileResp{}, nil
}

type userProfileData struct {
	User *entity.User
}

type persistingUserProfileData struct {
	UpdatingUser *entity.User
}

func (uc *UserUC) loadUserProfileData(
	ctx context.Context,
	db database.IDB,
	auth *basedto.Auth,
	req *userdto.UpdateProfileReq,
	data *userProfileData,
) error {
	user, err := uc.userRepo.GetByID(ctx, db, auth.User.ID,
		bunex.SelectFor("UPDATE"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}

	if user.Status != base.UserStatusActive {
		return apperrors.New(apperrors.ErrActionNotAllowed).
			WithMsgLog("user '%s' not active", user.Email)
	}
	data.User = user

	// Save user photo to local disk
	if req.Photo != nil && req.Photo.FileName != "" {
		err = uc.userService.SaveUserPhoto(ctx, user, req.Photo.DataBytes, req.Photo.FileExt)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func (uc *UserUC) preparePersistingUserProfileData(
	req *userdto.UpdateProfileReq,
	profileData *userProfileData,
	persistingData *persistingUserProfileData,
) {
	timeNow := timeutil.NowUTC()
	user := profileData.User
	persistingData.UpdatingUser = user

	user.UpdatedAt = timeNow
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Photo != nil && req.Photo.FileName == "" {
		user.Photo = ""
	}
}
