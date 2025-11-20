package secretuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/bunex"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
	"github.com/localpaas/localpaas/localpaas_app/pkg/transaction"
	"github.com/localpaas/localpaas/localpaas_app/usecase/providers/secretuc/secretdto"
)

func (uc *SecretUC) DeleteSecret(
	ctx context.Context,
	auth *basedto.Auth,
	req *secretdto.DeleteSecretReq,
) (*secretdto.DeleteSecretResp, error) {
	err := transaction.Execute(ctx, uc.db, func(db database.Tx) error {
		secretData := &deleteSecretData{}
		err := uc.loadSecretDataForDelete(ctx, db, req, secretData)
		if err != nil {
			return apperrors.Wrap(err)
		}

		persistingData := &persistingSecretData{}
		uc.prepareDeletingSecret(secretData, persistingData)

		return uc.persistData(ctx, db, persistingData)
	})
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	return &secretdto.DeleteSecretResp{}, nil
}

type deleteSecretData struct {
	Setting *entity.Setting
}

func (uc *SecretUC) loadSecretDataForDelete(
	ctx context.Context,
	db database.IDB,
	req *secretdto.DeleteSecretReq,
	data *deleteSecretData,
) error {
	options := []bunex.SelectQueryOption{
		bunex.SelectFor("UPDATE OF setting"),
	}
	if req.ObjectID != "" {
		options = append(options, bunex.SelectWhere("setting.object_id = ?", req.ObjectID))
	} else {
		options = append(options, bunex.SelectWhere("setting.object_id IS NULL"))
	}

	setting, err := uc.settingRepo.GetByID(ctx, db, base.SettingTypeSecret, req.ID, false, options...)
	if err != nil {
		return apperrors.Wrap(err)
	}
	data.Setting = setting

	return nil
}

func (uc *SecretUC) prepareDeletingSecret(
	data *deleteSecretData,
	persistingData *persistingSecretData,
) {
	setting := data.Setting
	setting.DeletedAt = timeutil.NowUTC()
	persistingData.UpsertingSettings = append(persistingData.UpsertingSettings, setting)
}
