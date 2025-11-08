package useruc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/useruc/userdto"
	"github.com/localpaas/localpaas/pkg/timeutil"
)

func (uc *UserUC) UpdateStatus(
	ctx context.Context,
	auth *basedto.Auth,
	req *userdto.UpdateStatusReq,
) (*userdto.UpdateStatusResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		profileData := &userStatusData{}
		err := uc.loadUserStatusData(ctx, db, req, profileData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingUserStatusData{}
		uc.preparePersistingUserStatusData(req, profileData, persistingData)

		return uc.persistUserStatusData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &userdto.UpdateStatusResp{}, nil
}

type userStatusData struct {
	User *entity.User
}

type persistingUserStatusData struct {
	UpdatingUser *entity.User
}

func (uc *UserUC) loadUserStatusData(
	ctx context.Context,
	db database.IDB,
	req *userdto.UpdateStatusReq,
	data *userStatusData,
) error {
	user, err := uc.userRepo.GetByID(ctx, db, req.ID,
		bunex.SelectFor("UPDATE"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.User = user

	return nil
}

func (uc *UserUC) preparePersistingUserStatusData(
	req *userdto.UpdateStatusReq,
	profileData *userStatusData,
	persistingData *persistingUserStatusData,
) {
	timeNow := timeutil.NowUTC()
	user := profileData.User
	persistingData.UpdatingUser = user

	user.UpdatedAt = timeNow
	user.Status = req.Status
}

func (uc *UserUC) persistUserStatusData(
	ctx context.Context,
	db database.IDB,
	persistingData *persistingUserStatusData,
) error {
	err := uc.userRepo.Update(ctx, db, persistingData.UpdatingUser,
		bunex.UpdateColumns("updated_at", "status"),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
